<?php

namespace App\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\Process;
use Tightenco\Collect\Support\Collection;

class LsCommand extends Command
{
    protected static $defaultName = 'app:ls';

    protected Collection $config;

    public function __construct(string $name = null)
    {

        $this->config = new Collection(json_decode(file_get_contents('config.json'), true)['instances']);
        parent::__construct($name);
    }

    protected function configure()
    {
        $this
            ->setDescription('Run ls command on all servers')
            ->setHelp('Run ls command on all servers');
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
