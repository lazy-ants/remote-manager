<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class UpgradeCommand extends AbstractCommand
{
    protected static $defaultName = 'app:upgrade';

    protected function configure()
    {
        $this
            ->setDescription('Upgrade server packages')
            ->setHelp('Upgrade server packages');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process(
                'echo $PASSWORD | sudo -S apt-get update' .
                '&& echo $PASSWORD | sudo -S apt-get -y upgrade' .
                '&& echo $PASSWORD | sudo -S apt-get -y autoremove'
            )
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
