<?php

namespace App\Command;

use Carbon\Carbon;
use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\InputStream;
use Symfony\Component\Process\Process;

class KernelCommand extends AbstractCommand
{
    protected static $defaultName = 'app:kernel';

    protected function configure()
    {
        $this
            ->setDescription('Get server kernels')
            ->setHelp('Get server kernels');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $progressBar = new ProgressBar($output, $this->config->count());
        $progressBar->start();

        $table = new Table($output);
        $table->setHeaders(['Name', 'Kernel']);

        $this->config->each(
            function ($hostConfig, $i) use ($output, $progressBar, $table) {
                $input = new InputStream();

                $process = new Process(['ssh', $hostConfig['connection-string']]);
                $process->setInput($input);
                $process->start();

                $input->write('echo "startoutputsisteminformation"');
                $input->write('&& uname -r');

                $input->close();

                $process->wait();

                $serverOutput = trim(explode('startoutputsisteminformation', $process->getOutput())[1]);

                $table->addRow(
                    [
                        $hostConfig['name'],
                        $serverOutput,
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
