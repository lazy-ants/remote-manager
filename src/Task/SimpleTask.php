<?php

namespace App\Task;

use Symfony\Component\Process\Process;

class SimpleTask extends AbstractTask
{
    public function run(): string
    {
        if ($this->needSudo) {
            if (empty($this->sudoPassword)) {
                throw new \InvalidArgumentException($this->name . ' needs sudo password');
            }
            $process = new Process(
                [
                    'ssh',
                    $this->connectionString,
                    '-o SendEnv="PASSWORD"',
                ],
                null,
                [
                    'PASSWORD' => $this->sudoPassword,
                ]
            );
        } else {
            $process = new Process(
                [
                    'ssh',
                    $this->connectionString,
                ]
            );
        }

        $process->setInput('echo "startoutputsysteminformation" && ' . $this->command);

        $process->run();

        $result = explode('startoutputsysteminformation', $process->getOutput());

        return trim(
            str_replace(
                ['\n', '\l'],
                '',
                isset($result[1]) ? $result[1] : ''
            )
        );
    }
}
