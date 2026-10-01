<?php

declare(strict_types=1);

namespace App\Support;

use DateTimeImmutable;
use DateTimeInterface;
use DateTimeZone;
use MongoDB\BSON\UTCDateTime;

final class Dates
{
    public static function now(): DateTimeImmutable
    {
        return new DateTimeImmutable('now', new DateTimeZone('UTC'));
    }

    public static function fromBson(UTCDateTime|DateTimeInterface $value): DateTimeImmutable
    {
        $dt = $value instanceof UTCDateTime ? $value->toDateTime() : $value;

        return DateTimeImmutable::createFromInterface($dt)->setTimezone(new DateTimeZone('UTC'));
    }

    public static function toIso(UTCDateTime|DateTimeInterface|null $value): ?string
    {
        if ($value === null) {
            return null;
        }

        return self::fromBson($value)->format(DateTimeInterface::RFC3339_EXTENDED);
    }
}
