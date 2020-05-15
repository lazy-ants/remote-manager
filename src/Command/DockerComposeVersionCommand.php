<?php

namespace App\Command;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class DockerComposeVersionCommand extends AbstractCommand
{
    protected static $defaultName = 'app:docker-compose-version';

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Get docket compose version')
            ->setHelp('Get docket compose version');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $table = new Table($output);
        $table->setHeaders(['Name', 'Docker compose version']);

        $this
            ->process('docker-compose -v')
            ->outputTable($output, $table)
            ->outputErrors($output);

        return 0;
    }
}
