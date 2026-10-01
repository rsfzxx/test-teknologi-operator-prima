<?php

declare(strict_types=1);

namespace App\Domain;

enum StoreStatus: string
{
    case Active = 'Active';
    case Inactive = 'Inactive';

    public static function fromInput(string $value): ?self
    {
        return match (strtolower(trim($value))) {
            'active'   => self::Active,
            'inactive' => self::Inactive,
            default    => null,
        };
    }

    /** @return list<string> */
    public static function values(): array
    {
        return array_map(static fn (self $s): string => $s->value, self::cases());
    }
}
