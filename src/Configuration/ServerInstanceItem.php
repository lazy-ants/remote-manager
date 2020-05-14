<?php

namespace App\Configuration;

class ServerInstanceItem
{
    public string $name = '';
    public string $connectionString = '';
    public string $sudoPassword = '';
    public array $tags = [];
}
