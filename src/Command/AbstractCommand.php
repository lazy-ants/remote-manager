<?php

namespace App\Command;

use App\Configuration\ServerInstanceItem;
use App\Configuration\ServerInstancesConfig;
use App\Task\AbstractTask;
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
    protected bool $needSudo = false;

    public function __construct(string $name = null)
    {
        $this->results = new Collection();
        $this->errors = new Collection();
        $this->timeouts = new Collection();

        parent::__construct($name);
    }

    protected function init()
    {
        $this->pool = Pool::create();
        $this->config = new ServerInstancesConfig();

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
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $this->init();

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
        /** @var ServerInstanceItem $hostConfig */
        foreach ($this->config as $i => $hostConfig) {
            $this->pool
                ->add(
                    $this->needSudo ?
                        new SimpleTask(
                            $hostConfig->name,
                            $hostConfig->connectionString,
                            $command,
                            AbstractTask::NEED_SUDO,
                            $hostConfig->sudoPassword
                        ) :
                        new SimpleTask(
                            $hostConfig->name,
                            $hostConfig->connectionString,
                            $command
                        ),
                    1024 * 100
                )
                ->then(
                    function ($output) use ($hostConfig) {
                        $this->results->push(
                            [
                                'name' => $hostConfig->name,
                                'value' => $output,
                            ]
                        );

                        $this->progressBar->advance();
                    }
                )
                ->catch(
                    function (Throwable $exception) use ($hostConfig) {
                        $this->errors->push(
                            [
                                'hostName' => $hostConfig->name,
                                'code' => $exception->getCode(),
                                'file' => $exception->getFile(),
                                'line' => $exception->getLine(),
                                'message' => $exception->getMessage(),
                            ]
                        );
                    }
                )
                ->timeout(
                    function () use ($hostConfig) {
                        $this->timeouts->push($hostConfig->name);
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
                function ($item, $i) use ($output) {
                    if (0 == $i) {
                        $output->writeln('<error>Errors</error>');
                    }
                    $output->writeln('host name: ' . $item['hostName']);
                    $output->writeln('code: ' . $item['code']);
                    $output->writeln('file: ' . $item['file']);
                    $output->writeln('line: ' . $item['line']);
                    $output->writeln('message: ' . $item['message']);
                }
            );

        $this->timeouts
            ->each(
                function ($item, $i) use ($output) {
                    if (0 == $i) {
                        $output->writeln('<error>Timeouts</error>');
                    }
                    $output->writeln($item);
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
                function ($item, $i) use ($output, $total) {
                    $output->writeln(sprintf('<info>[%s / %s] %s:</info>', $i + 1, $total, $item['name']));
                    $output->writeln($item['value']);
                }
            );

        return $this;
    }
}
