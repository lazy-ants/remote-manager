<?php

namespace App\Command;

use Symfony\Component\Console\Helper\ProgressBar;
use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Tightenco\Collect\Support\Collection;

class CheckRebootCommand extends AbstractCommand
{
    protected static $defaultName = 'app:check-reboot';
    protected bool $needSudo = true;

    protected function configure()
    {
        parent::configure();

        $this
            ->setDescription('Checks whether a reboot is required')
            ->setHelp('Checks whether a reboot is required')
            ->addOption('reboot', 'r', InputOption::VALUE_NONE, 'Reboot the server if needed');
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

        if ($input->getOption('reboot')) {
            $backupConfig = clone $this->config;

            $rebootInstances = $this->results
                ->where('value', 'yes')
                ->except('value')
                ->flatten()
                ->all();
            $this->config->filterByNames($rebootInstances);

            if ($this->config->instances->count() > 0) {
                $table = new Table($output);
                $table->setHeaders(['Name', 'Reboot started']);

                $output->writeln('<info>Total servers to reboot:</info> ' . $this->config->instances->count());

                $this->progressBar = new ProgressBar($output, $this->config->instances->count());
                $this->progressBar->start();

                $this->results = new Collection();
                $this->process('echo $PASSWORD | sudo -S reboot');

                $this->results->transform(
                    function ($item) {
                        $item['value'] = 'yes';

                        return $item;
                    }
                );

                $this
                    ->outputTable($output, $table)
                    ->outputErrors($output);
            }

            $this->config = $backupConfig;
        }

        return 0;
    }
}
