<?php

declare(strict_types=1);

namespace App\Database;

use Closure;
use MongoDB\BSON\UTCDateTime;
use MongoDB\Database;
use RuntimeException;
use Throwable;

final class Migrator
{
    private const COLLECTION = 'migrations';

    public function __construct(
        private readonly Database $database,
        private readonly string $directory,
    ) {
    }

    /**
     * @param (Closure(string): void)|null $log
     * @return list<string> nama migrasi yang dijalankan
     */
    public function run(?Closure $log = null): array
    {
        $log ??= static function (string $message): void {
        };

        $history = $this->database->selectCollection(self::COLLECTION);
        $history->createIndex(['name' => 1], ['unique' => true, 'name' => 'uq_migrations_name']);

        $files = glob($this->directory . DIRECTORY_SEPARATOR . '*.php') ?: [];
        sort($files, SORT_STRING);

        $applied = [];
        foreach ($files as $file) {
            $name = basename($file);

            if ($history->countDocuments(['name' => $name], ['limit' => 1]) > 0) {
                continue;
            }

            $migration = require $file;
            if (!$migration instanceof Closure) {
                throw new RuntimeException(sprintf('Migration %s must return a Closure', $name));
            }

            try {
                $migration($this->database);
            } catch (Throwable $e) {
                throw new RuntimeException(sprintf('Migration %s failed: %s', $name, $e->getMessage()), 0, $e);
            }

            $history->insertOne([
                'name'       => $name,
                'applied_at' => new UTCDateTime(),
            ]);

            $log('migrasi dijalankan: ' . $name);
            $applied[] = $name;
        }

        return $applied;
    }
}
