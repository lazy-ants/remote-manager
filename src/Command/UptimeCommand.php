<?php

namespace App\Command;

use Carbon\Carbon;
use Spatie\Async\Pool;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\Process;
use Throwable;
use Tightenco\Collect\Support\Collection;

class UptimeCommand extends AbstractCommand
{
    protected static $defaultName = 'app:uptime';

    protected function configure()
    {
        $this
            ->setDescription('Get server uptime')
            ->setHelp('Get server uptime');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $progressBar = new ProgressBar($output, $this->config->count());
        $progressBar->start();

        $pool = Pool::create();
        $results = new Collection();

        foreach ($this->config as $i => $hostConfig) {
            $pool
                ->add(
                    function () use ($hostConfig) {
                        $process = new Process(['ssh', $hostConfig['connection-string']]);
                        $process->setInput('echo "startoutputsisteminformation" && uptime -s');

                        $process->run();

                        $serverOutput = trim(
                            str_replace(
                                ['\n', '\l'],
                                '',
                                explode('startoutputsisteminformation', $process->getOutput())[1]
                            )
                        );

                        return Carbon::createFromFormat('Y-m-d H:i:s', $serverOutput)->diffInDays() . ' days';
                    }
                )
                ->then(
                    function ($output) use (&$results, $hostConfig, $progressBar) {
                        $results->push(
                            [
                                'name' => $hostConfig['name'],
                                'value' => $output,
                            ]
                        );
                        $progressBar->advance();
                    }
                )
                ->catch(
                    function (Throwable $exception) {
                        var_dump($exception);
                    }
                )
                ->timeout(
                    function () {
                        var_dump('timeout');
                    }
                );
        }

        $pool->wait();

        $progressBar->finish();
        $output->writeln('');

        $table = new Table($output);
        $table->setHeaders(['Name', 'OS']);

        $results
            ->sortBy('name')
            ->each(
                function ($body) use ($table) {
                    $table->addRow(
                        [
                            $body['name'],
                            $body['value'],
                        ]
                    );

                }
            );

        $table->render();

        return 0;
    }
}
