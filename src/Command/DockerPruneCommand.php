<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class DockerPruneCommand extends AbstractCommand
{
    protected static $defaultName = 'app:docker-prune';

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Prune old docker data')
            ->setHelp('Prune old docker data');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process('echo "y" | docker system prune')
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
