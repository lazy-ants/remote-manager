<?php

namespace App\Configuration;

use Tightenco\Collect\Support\Collection;

class ServerInstancesConfig
{
    public Collection $instances;

    public function __construct()
    {
        $this->instances = new Collection();

        foreach (json_decode(file_get_contents('config.json'), true)['instances'] as $item) {
            $serverInstanceItem = new ServerInstanceItem();

            $serverInstanceItem->name = $item['name'];
            $serverInstanceItem->connectionString = $item['connection-string'];
            $serverInstanceItem->sudoPassword = !empty($item['sudo-password']) ?
                $item['sudo-password'] :
                $_ENV['SUDO_PASSWORD'];
            $serverInstanceItem->tags = explode(',', !empty($item['tags']) ? $item['tags'] : '');

            $this->instances->push($serverInstanceItem);
        }
    }

    public function filterByTags(array $tags = []): ServerInstancesConfig
    {
        $this->instances = $this->instances
            ->filter(
                function (ServerInstanceItem $item) use ($tags) {
                    return collect($tags)
                            ->filter()
                            ->diff($item->tags)
                            ->count() == 0;
                }
            )
            ->values();

        return $this;
    }
}
