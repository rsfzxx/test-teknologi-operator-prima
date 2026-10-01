<?php

declare(strict_types=1);

use MongoDB\Database;

return static function (Database $db): void {
    $validator = [
        '$jsonSchema' => [
            'bsonType' => 'object',
            'required' => ['store_uuid', 'name', 'address', 'email', 'phone', 'status', 'created_at', 'updated_at'],
            'properties' => [
                'store_uuid' => [
                    'bsonType'    => 'string',
                    'pattern'     => '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$',
                    'description' => 'UUID v4 dalam huruf kecil',
                ],
                'name' => [
                    'bsonType'  => 'string',
                    'minLength' => 1,
                    'maxLength' => 150,
                ],
                'address' => [
                    'bsonType'  => 'string',
                    'minLength' => 1,
                    'maxLength' => 255,
                ],
                'email' => [
                    'bsonType'  => 'string',
                    'maxLength' => 254,
                    'pattern'   => '^[^@\s]+@[^@\s]+\.[^@\s]+$',
                ],
                'phone' => [
                    'bsonType' => 'string',
                    'pattern'  => '^\+?[0-9]{8,15}$',
                ],
                'status' => [
                    'enum' => ['Active', 'Inactive'],
                ],
                'notes' => [
                    'bsonType'  => ['string', 'null'],
                    'maxLength' => 500,
                ],
                'created_at' => ['bsonType' => 'date'],
                'updated_at' => ['bsonType' => 'date'],
            ],
        ],
    ];

    $exists = iterator_to_array($db->listCollectionNames(['filter' => ['name' => 'stores']]));

    if ($exists === []) {
        $db->createCollection('stores', [
            'validator'        => $validator,
            'validationLevel'  => 'strict',
            'validationAction' => 'error',
        ]);
        return;
    }

    $db->command([
        'collMod'          => 'stores',
        'validator'        => $validator,
        'validationLevel'  => 'strict',
        'validationAction' => 'error',
    ]);
};
