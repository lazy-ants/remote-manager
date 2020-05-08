<?php

namespace App\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\Process;
use Tightenco\Collect\Support\Collection;

class UpgradeCommand extends Command
{
    protected static $defaultName = 'app:upgrade';

    protected Collection $config;

    public function __construct(string $name = null)
    {

        $this->config = new Collection(json_decode(file_get_contents('config.json'), true)['instances']);
        parent::__construct($name);
    }

    protected function configure()
    {
        $this
            ->setDescription('Upgrade server packages')
            ->setHelp('Upgrade server packages');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
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

        $output->writeln('<info>Total servers:</info> ' . $this->config->count());

        $this->config->each(
            function ($hostConfig, $i) use ($output) {
                $output->writeln('');
                $output->writeln(sprintf('<info>[%s] Running:</info> %s', $i + 1, $hostConfig['name']));

                $process = new Process(
                    [
                        'ssh',
                        $hostConfig['connection-string'],
                        '-o SendEnv="PASSWORD"',
                        'echo $PASSWORD | sudo -S apt-get update',
                    ],
                    null,
                    [
                        'PASSWORD' => $_ENV['SUDO_PASSWORD'],
                    ]
                );
                $process->setTty(Process::isTtySupported());
                $process->start();
                $process->wait();

                $process = new Process(
                    [
                        'ssh',
                        $hostConfig['connection-string'],
                        '-o SendEnv="PASSWORD"',
                        'echo $PASSWORD | sudo -S apt-get -y upgrade',
                    ],
                    null,
                    [
                        'PASSWORD' => $_ENV['SUDO_PASSWORD'],
                    ]
                );
                $process->setTty(Process::isTtySupported());
                $process->start();
                $process->wait();

                $process = new Process(
                    [
                        'ssh',
                        $hostConfig['connection-string'],
                        '-o SendEnv="PASSWORD"',
                        'echo $PASSWORD | sudo -S apt-get -y autoremove',
                    ],
                    null,
                    [
                        'PASSWORD' => $_ENV['SUDO_PASSWORD'],
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
