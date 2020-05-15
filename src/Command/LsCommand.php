<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class LsCommand extends AbstractCommand
{
    protected static $defaultName = 'app:ls';

    protected function configure()
    {
        $this
            ->setDescription('Run ls command on all servers')
            ->setHelp('Run ls command on all servers')
            ->addArgument('arg', InputArgument::OPTIONAL, 'ls command arguments', '-lha');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process('ls ' . $input->getArgument('arg'))
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
