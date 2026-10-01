<?php

declare(strict_types=1);

namespace App\Database;

use MongoDB\Client;

final class MongoConnection
{
    private const TYPE_MAP = [
        'root'     => 'array',
        'document' => 'array',
        'array'    => 'array',
    ];

    public static function createClient(string $uri): Client
    {
        return new Client(
            $uri,
            [
                'appName'                  => 'operator-prima-slim',
                'serverSelectionTimeoutMS' => 5000,
            ],
            ['typeMap' => self::TYPE_MAP],
        );
    }
}
