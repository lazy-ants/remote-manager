<?php

namespace App\Command;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class CheckRebootCommand extends AbstractCommand
{
    protected static $defaultName = 'app:check-reboot';

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Checks whether a reboot is required')
            ->setHelp('Checks whether a reboot is required');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this->process('ls /var/run/reboot-required');

        $table = new Table($output);
        $table->setHeaders(['Name', 'Reboot required?']);

        $this->results->transform(
            function ($item) {
                $item['value'] = '/var/run/reboot-required' == $item['value'] ? 'yes' : 'no';

                return $item;
            }
        );

        $this
            ->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
