<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\Process;

class LsCommand extends AbstractCommand
{
    protected static $defaultName = 'app:ls';

    protected function configure()
    {
        $this
            ->setDescription('Run ls command on all servers')
            ->setHelp('Run ls command on all servers');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $this->config->each(
            function ($hostConfig, $i) use ($output) {
                $output->writeln('');
                $output->writeln(sprintf('<info>[%s] Running:</info> %s', $i + 1, $hostConfig['name']));

                $process = new Process(
                    [
                        'ssh',
                        $hostConfig['connection-string'],
                        'ls -lha',
                    ]
                );
                $process->setTty(Process::isTtySupported());
                $process->start();
                $process->wait();

                $output->writeln('<info>Finished:</info> ' . $hostConfig['name']);
            }
        );

        return 0;
    }
}
