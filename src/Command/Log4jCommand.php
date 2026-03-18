<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class Log4jCommand extends AbstractCommand
{
    protected static $defaultName = 'app:log4j';
    protected bool $needSudo = true;

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Check if Log4j is used on the server')
            ->setHelp('Run the check on all servers');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process('echo $PASSWORD | sudo -S ls -lha')
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
