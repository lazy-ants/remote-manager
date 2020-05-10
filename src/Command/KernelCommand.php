<?php

namespace App\Command;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

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
        parent::execute($input, $output);

        $table = new Table($output);
        $table->setHeaders(['Name', 'Kernel']);

        $this
            ->process('uname -r')
            ->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
