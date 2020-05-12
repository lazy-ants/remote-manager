<?php

namespace App\Task;

abstract class AbstractTask
{
    const NEED_SUDO = true;

    protected string $name;
    protected string $connectionString;
    protected string $sudoPassword;
    protected string $command;
    protected bool $needSudo;

    public function __construct(
        string $name,
        string $connectionString,
        string $command,
        bool $needSudo = false,
        string $sudoPassword = ''
    ) {
        $this->name = $name;
        $this->connectionString = $connectionString;
        $this->command = $command;
        $this->needSudo = $needSudo;
        $this->sudoPassword = $sudoPassword;
    }

    public function __invoke()
    {
        return $this->run();
    }

    abstract public function run();
}
