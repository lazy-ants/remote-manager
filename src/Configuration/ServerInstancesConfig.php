<?php

namespace App\Configuration;

use Tightenco\Collect\Support\Collection;

class ServerInstancesConfig extends Collection
{
    public function __construct()
    {
        foreach (json_decode(file_get_contents('config.json'), true)['instances'] as $item) {
            $serverInstanceItem = new ServerInstanceItem();

            $serverInstanceItem->name = $item['name'];
            $serverInstanceItem->connectionString = $item['connection-string'];
            $serverInstanceItem->sudoPassword = !empty($item['sudo-password']) ?
                $item['sudo-password'] :
                $_ENV['SUDO_PASSWORD'];

            $this->push($serverInstanceItem);
        }
    }
}
