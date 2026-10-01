<?php

declare(strict_types=1);

use MongoDB\Database;

return static function (Database $db): void {
    $db->selectCollection('stores')->createIndexes([
        ['key' => ['store_uuid' => 1], 'name' => 'uq_stores_store_uuid', 'unique' => true],
        ['key' => ['email' => 1], 'name' => 'uq_stores_email', 'unique' => true],
        ['key' => ['status' => 1], 'name' => 'idx_stores_status'],
        ['key' => ['created_at' => -1], 'name' => 'idx_stores_created_at'],
    ]);
};
