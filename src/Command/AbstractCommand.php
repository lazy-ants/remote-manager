<?php

namespace App\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Process\Process;
use Tightenco\Collect\Support\Collection;

abstract class AbstractCommand extends Command
{
    protected Collection $config;

    public function __construct(string $name = null)
    {
        $this->config = new Collection(json_decode(file_get_contents('config.json'), true)['instances']);

        # add private keys to the ssh agent
        if (!empty($_ENV['PPK_NAMES'])) {
            foreach (explode(',', $_ENV['PPK_NAMES']) as $ppkName) {
                $process = new Process(
                    [
                        'ssh-add',
                        '/root/.ssh/' . $ppkName,
                    ]
                );
                $process->setTty(Process::isTtySupported());
                $process->start();
                $process->wait();
            }
        }

        parent::__construct($name);
    }
}
