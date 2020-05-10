<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

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
        parent::execute($input, $output);

        $this
            ->process('ls -lha')
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
