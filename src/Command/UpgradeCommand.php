<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class UpgradeCommand extends AbstractCommand
{
    protected static $defaultName = 'app:upgrade';
    protected bool $needSudo = true;

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Upgrade server packages [need sudo]')
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
