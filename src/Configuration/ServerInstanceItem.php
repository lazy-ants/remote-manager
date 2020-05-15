<?php

namespace App\Configuration;

class ServerInstanceItem
{
    public string $name = '';
    public string $connectionString = '';
    public string $sudoPassword = '';

    /**
     * @var array<string>
     */
    public array $tags = [];
}
