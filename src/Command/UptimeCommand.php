<?php

namespace App\Command;

use Carbon\Carbon;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\InputStream;
use Symfony\Component\Process\Process;

class UptimeCommand extends AbstractCommand
{
    protected static $defaultName = 'app:uptime';

    protected function configure()
    {
        $this
            ->setDescription('Get server uptimes')
            ->setHelp('Get server uptimes');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $progressBar = new ProgressBar($output, $this->config->count());
        $progressBar->start();

        $table = new Table($output);
        $table->setHeaders(['Name', 'Uptime']);

        $this->config->each(
            function ($hostConfig, $i) use ($output, $progressBar, $table) {
                $input = new InputStream();

                $process = new Process(['ssh', $hostConfig['connection-string']]);
                $process->setInput($input);
                $process->start();

                $input->write('echo "startoutputsisteminformation"');
                $input->write('&& uptime -s');

                $input->close();

                $process->wait();

                $serverOutput = trim(explode('startoutputsisteminformation', $process->getOutput())[1]);

                $table->addRow(
                    [
                        $hostConfig['name'],
                        Carbon::createFromFormat('Y-m-d H:i:s', $serverOutput)->diffInDays() . ' days',
                    ]
                );


                $progressBar->advance();
            }
        );

        $progressBar->finish();
        $output->writeln('');
        $table->render();

        return 0;
    }
}
