<?php

namespace App\Command;

use App\Configuration\ServerInstanceItem;
use App\Task\AbstractTask;
use App\Task\SimpleTask;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Process\InputStream;

class ValidateConfigCommand extends AbstractCommand
{
    protected static $defaultName = 'app:validate-config';
    protected bool $needSudo = true;

    protected function configure()
    {
        $this
            ->setDescription('Validate server instances config')
            ->setHelp('Validate server instances config');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $this->init();

        $table = new Table($output);
        $table->setHeaders(['Name', 'Login possible?', 'Sudo password exposed?', 'Sudo possible?']);

        $this->config->each(
            function (ServerInstanceItem $item, $i) use ($output, $table) {
                $tableRow = [$item->name];
                $output->writeln(
                    sprintf(
                        'Checking %s of %s: %s',
                        $i + 1,
                        $this->config->count(),
                        $item->name
                    )
                );

                $input = new InputStream();

                # check if login possible
                $task = new SimpleTask($item->name, $item->connectionString, 'whoami');
                $result = $task->run();
                $tableRow[] = !empty($result) ? 'yes' : 'no';

                # check if sudo password exposed
                if (!empty($item->sudoPassword)) {
                    $task = new SimpleTask(
                        $item->name,
                        $item->connectionString,
                        '[ -z "$PASSWORD" ] && echo "no" || echo "yes"',
                        AbstractTask::NEED_SUDO,
                        $item->sudoPassword
                    );
                    $result = $task->run();
                    $tableRow[] = $result;

                    # check if sudo possible
                    if ('yes' == $result) {
                        $task = new SimpleTask(
                            $item->name,
                            $item->connectionString,
                            'echo $PASSWORD | sudo -S whoami',
                            AbstractTask::NEED_SUDO,
                            $item->sudoPassword
                        );
                        $result = $task->run();

                        $tableRow[] = !empty($result) ? 'yes' : 'no';
                    } else {
                        $tableRow[] = 'no';
                    }
                }

                $table->addRow($tableRow);
            }
        );

        $output->writeln('');
        $table->render();

        return 0;
    }
}
