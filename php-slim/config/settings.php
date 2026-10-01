<?php

declare(strict_types=1);

return [
    'app' => [
        'debug' => filter_var($_ENV['APP_DEBUG'] ?? 'false', FILTER_VALIDATE_BOOL),
    ],
    'mongodb' => [
        'uri'      => $_ENV['MONGODB_URI'] ?? '',
        'database' => $_ENV['MONGODB_DATABASE'] ?? 'operator_prima',
    ],
    'paths' => [
        'migrations' => dirname(__DIR__) . '/database/migrations',
        'pipelines'  => dirname(__DIR__) . '/pipelines',
    ],
];
