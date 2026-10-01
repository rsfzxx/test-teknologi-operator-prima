<?php

declare(strict_types=1);

namespace App\Database;

use JsonException;
use RuntimeException;

final class PipelineLoader
{
    /** @var array<string, array> */
    private array $cache = [];

    public function __construct(private readonly string $directory)
    {
    }

    public function load(string $name): array
    {
        if (isset($this->cache[$name])) {
            return $this->cache[$name];
        }

        $path = $this->directory . DIRECTORY_SEPARATOR . $name . '.json';
        if (!is_file($path)) {
            throw new RuntimeException(sprintf('Pipeline file not found: %s', $path));
        }

        try {
            $pipeline = json_decode((string) file_get_contents($path), true, 512, JSON_THROW_ON_ERROR);
        } catch (JsonException $e) {
            throw new RuntimeException(sprintf('Pipeline %s is not valid JSON: %s', $name, $e->getMessage()), 0, $e);
        }

        if (!is_array($pipeline) || !array_is_list($pipeline)) {
            throw new RuntimeException(sprintf('Pipeline %s must be a JSON array of stages', $name));
        }

        return $this->cache[$name] = $pipeline;
    }
}
