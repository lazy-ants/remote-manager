<?php

namespace App\Command;

use Carbon\Carbon;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class UptimeCommand extends AbstractCommand
{
    protected static $defaultName = 'app:uptime';

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Get server uptime')
            ->setHelp('Get server uptime');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $table = new Table($output);
        $table->setHeaders(['Name', 'Uptime']);

        $this->process('uptime -s');

        $this->results->transform(
            function ($item) {
                $item['value'] = Carbon::createFromFormat('Y-m-d H:i:s', $item['value'])->diffInDays() . ' days';

                return $item;
            }
        );

        $this
            ->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
