<?php

namespace App\Command;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class UfwCommand extends AbstractCommand
{
    protected static $defaultName = 'app:ufw';
    protected bool $needSudo = true;

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Get ufw status [need sudo]')
            ->setHelp('Get ufw status [need sudo]');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $table = new Table($output);
        $table->setHeaders(['Name', 'UFW Status']);

        $this->process('echo $PASSWORD | sudo -S ufw status');

        $this->results->transform(
            function ($item) {
                preg_match("/Status: (.*)\n/m", $item['value'], $matches);

                $item['value'] = isset($matches[1]) ? $matches[1] : '';

                return $item;
            }
        );

        $this->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
