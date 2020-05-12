<?php

namespace App\Command;

use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class SystemInfoCommand extends AbstractCommand
{
    protected static $defaultName = 'app:system-info';
    protected bool $needSudo = true;

    protected function configure()
    {
        $this
            ->setDescription('Get system information')
            ->setHelp('Get system information');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        parent::execute($input, $output);

        $this
            ->process(
                sprintf(
                    'echo "%s" > /tmp/system-info.base64',
                    base64_encode(file_get_contents('src/Scripts/system-info.sh'))
                ) .
                '&& base64 -d /tmp/system-info.base64 > /tmp/system-info.sh' .
                '&& rm /tmp/system-info.base64' .
                '&& chmod +x /tmp/system-info.sh' .
                '&& echo $PASSWORD | sudo -S /tmp/system-info.sh' .
                '&& rm /tmp/system-info.sh'
            )
            ->outputList($output)
            ->outputErrors($output);

        return 0;
    }
}
