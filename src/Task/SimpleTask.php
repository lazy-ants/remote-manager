<?php

namespace App\Task;

use Symfony\Component\Process\Process;

class SimpleTask extends AbstractTask
{
    public function run()
    {
        $process = new Process(
            [
                'ssh',
                $this->connectionString,
                '-o SendEnv="PASSWORD"',
            ],
            null,
            [
                'PASSWORD' => $_ENV['SUDO_PASSWORD'],
            ]
        );

        $process->setInput('echo "startoutputsysteminformation" && ' . $this->command);

        $process->run();

        return trim(
            str_replace(
                ['\n', '\l'],
                '',
                explode('startoutputsysteminformation', $process->getOutput())[1]
            )
        );
    }
}