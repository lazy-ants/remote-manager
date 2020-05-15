<?php

namespace App\Command;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class OSCommand extends AbstractCommand
{
    protected static $defaultName = 'app:os';

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Get server OS')
            ->setHelp('Get server OS');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $table = new Table($output);
        $table->setHeaders(['Name', 'OS']);

        $this
            ->process('cat /etc/issue')
            ->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
