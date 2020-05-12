<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class DockerPsCommand extends AbstractCommand
{
    protected static $defaultName = 'app:docker-ps';

    protected function configure()
    {
        $this
            ->setDescription('Show docker process status')
            ->setHelp('Show docker process status');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process('docker ps')
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
