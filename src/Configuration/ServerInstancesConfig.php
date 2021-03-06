<?php

namespace App\Configuration;

use LogicException;
use Tightenco\Collect\Support\Collection;

class ServerInstancesConfig
{
    /**
     * @var Collection<ServerInstanceItem>
     */
    public Collection $instances;

    public function __construct()
    {
        $this->instances = new Collection();

        foreach (json_decode((string)file_get_contents('config.json'), true)['instances'] as $item) {
            $serverInstanceItem = new ServerInstanceItem();

            $serverInstanceItem->name = $item['name'];
            $serverInstanceItem->connectionString = $item['connection-string'];
            $serverInstanceItem->sudoPassword = !empty($item['sudo-password']) ?
                $item['sudo-password'] :
                $_ENV['SUDO_PASSWORD'];
            $serverInstanceItem->tags = explode(',', !empty($item['tags']) ? $item['tags'] : '');

            $this->instances->push($serverInstanceItem);
        }

        $duplicates = $this->instances->duplicates('name');

        if ($duplicates->count() > 0) {
            throw new LogicException(
                'Duplicate server names in the config.json found: ' . $duplicates->unique()->implode(', ')
            );
        }
    }

    /**
     * @param array<string> $tags
     * @return $this
     */
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

    /**
     * @param array<string> $names
     * @return $this
     */
    public function filterByNames(array $names = []): ServerInstancesConfig
    {
        $this->instances = $this->instances
            ->filter(
                function (ServerInstanceItem $item) use ($names) {
                    return in_array($item->name, $names);
                }
            )
            ->values();

        return $this;
    }

}
