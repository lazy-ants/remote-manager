<?php

namespace App\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\InputStream;
use Symfony\Component\Process\Process;
use Tightenco\Collect\Support\Collection;

class SystemInfoCommand extends Command
{
    protected static $defaultName = 'app:system-info';

    protected Collection $config;

    public function __construct(string $name = null)
    {

        $this->config = new Collection(json_decode(file_get_contents('config.json'), true)['instances']);
        parent::__construct($name);
    }

    protected function configure()
    {
        $this
            ->setDescription('Get system information')
            ->setHelp('Get system information');
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

                $input = new InputStream();

                $process = new Process(
                    [
                        'ssh',
                        $hostConfig['connection-string'],
                        '-o SendEnv="PASSWORD"',
                    ],
                    null,
                    [
                        'PASSWORD' => $_ENV['SUDO_PASSWORD'],
                    ]
                );
                $process->setInput($input);
                $process->start();

                $input->write('echo "' . base64_encode(file_get_contents('src/Scripts/system-info.sh')) . '" > /tmp/system-info.base64');
                $input->write('&& base64 -d /tmp/system-info.base64 > /tmp/system-info.sh');
                $input->write('&& rm /tmp/system-info.base64');
                $input->write('&& chmod +x /tmp/system-info.sh');
                $input->write('&& echo "startoutputsisteminformation"');
                $input->write('&& echo $PASSWORD | sudo -S /tmp/system-info.sh');
                $input->write('&& rm /tmp/system-info.sh');

                $input->close();

                $process->wait();

                $serverOutput = explode( 'startoutputsisteminformation',$process->getOutput())[1];

                $output->writeln($serverOutput);
                $output->writeln('<info>Finished:</info> ' . $hostConfig['name']);
            }
        );

        return 0;
    }
}
