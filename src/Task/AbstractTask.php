<?php

namespace App\Task;

abstract class AbstractTask
{
    protected string $connectionString;
    protected string $sudoPassword;
    protected string $command;

    public function __construct(string $connectionString, string $sudoPassword, string $command)
    {
        $this->connectionString = $connectionString;
        $this->sudoPassword = $sudoPassword;
        $this->command = $command;
    }

    public function __invoke()
    {
        return $this->run();
    }

    abstract public function run();
}
