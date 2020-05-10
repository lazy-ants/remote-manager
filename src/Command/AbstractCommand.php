<?php

namespace App\Command;

use App\Task\SimpleTask;
use Spatie\Async\Pool;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\Process;
use Throwable;
use Tightenco\Collect\Support\Collection;

abstract class AbstractCommand extends Command
{
    protected Collection $config;
    protected ProgressBar $progressBar;
    protected Collection $results;
    protected Pool $pool;
    protected Collection $errors;
    protected Collection $timeouts;

    public function __construct(string $name = null)
    {
        $this->config = new Collection(json_decode(file_get_contents('config.json'), true)['instances']);

        # add private keys to the ssh agent
        if (!empty($_ENV['PPK_NAMES'])) {
            foreach (explode(',', $_ENV['PPK_NAMES']) as $ppkName) {
                $process = new Process(
                    [
                        'ssh-add',
                        '/root/.ssh/' . $ppkName,
                    ]
                );
                $process->setTty(Process::isTtySupported());
                $process->start();
                $process->wait();
            }
        }

        $this->results = new Collection();
        $this->pool = Pool::create();
        $this->errors = new Collection();
        $this->timeouts = new Collection();

        parent::__construct($name);
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $this->progressBar = new ProgressBar($output, $this->config->count());
        $this->progressBar->start();
    }

    /**
     * @param string $command
     * @return $this
     */
    protected function process(string $command)
    {
        foreach ($this->config as $i => $hostConfig) {
            $hostName = $hostConfig['name'];

            $this->pool
                ->add(
                    new SimpleTask($hostConfig['connection-string'], $_ENV['SUDO_PASSWORD'], $command),
                    1024 * 100
                )
                ->then(
                    function ($output) use ($hostConfig) {
                        $this->results->push(
                            [
                                'name' => $hostConfig['name'],
                                'value' => $output,
                            ]
                        );

                        $this->progressBar->advance();
                    }
                )
                ->catch(
                    function (Throwable $exception) use ($hostName) {
                        $this->errors->push(
                            [
                                'hostName' => $hostName,
                                'code' => $exception->getCode(),
                                'file' => $exception->getFile(),
                                'line' => $exception->getLine(),
                                'message' => $exception->getMessage(),
                            ]
                        );
                    }
                )
                ->timeout(
                    function () use ($hostName) {
                        $this->timeouts->push($hostName);
                    }
                );
        }

        $this->pool->wait();

        $this->progressBar->finish();

        return $this;
    }

    /**
     * @param OutputInterface $output
     * @return $this
     */
    protected function outputErrors(OutputInterface $output)
    {
        $this->errors
            ->each(
                function ($error, $i) use ($output) {
                    if (0 == $i) {
                        $output->writeln('<error>Errors</error>');
                    }
                    $output->writeln('host name: ' . $error['hostName']);
                    $output->writeln('code: ' . $error['code']);
                    $output->writeln('file: ' . $error['file']);
                    $output->writeln('line: ' . $error['line']);
                    $output->writeln('message: ' . $error['message']);
                }
            );

        $this->timeouts
            ->each(
                function ($hostname, $i) use ($output) {
                    if (0 == $i) {
                        $output->writeln('<error>Timeouts</error>');
                    }
                    $output->writeln($hostname);
                }
            );

        return $this;
    }

    /**
     * @param OutputInterface $output
     * @param Table $table
     * @return $this
     */
    protected function outputTable(OutputInterface $output, Table $table)
    {
        $output->writeln('');

        $this->results
            ->sortBy('name')
            ->each(
                function ($item) use ($table) {
                    $table->addRow(
                        [
                            $item['name'],
                            $item['value'],
                        ]
                    );

                }
            );

        $table->render();

        return $this;
    }

    /**
     * @param OutputInterface $output
     * @return $this
     */
    protected function outputList(OutputInterface $output)
    {
        $output->writeln('');

        $total = $this->results->count();
        $this->results
            ->sortBy('name')
            ->values()
            ->each(
                function ($body, $i) use ($output, $total) {
                    $output->writeln(sprintf('<info>[%s / %s] %s:</info>', $i + 1, $total, $body['name']));
                    $output->writeln($body['value']);
                }
            );

        return $this;
    }
}
