<?php

declare(strict_types=1);

namespace App\Domain;

use App\Support\Dates;
use DateTimeImmutable;
use MongoDB\BSON\UTCDateTime;

final class Store
{
    public function __construct(
        public readonly string $storeUuid,
        public readonly string $name,
        public readonly string $address,
        public readonly string $email,
        public readonly string $phone,
        public readonly StoreStatus $status,
        public readonly ?string $notes,
        public readonly DateTimeImmutable $createdAt,
        public readonly DateTimeImmutable $updatedAt,
    ) {
    }

    public static function fromDocument(array $doc): self
    {
        return new self(
            storeUuid: (string) $doc['store_uuid'],
            name: (string) $doc['name'],
            address: (string) $doc['address'],
            email: (string) $doc['email'],
            phone: (string) $doc['phone'],
            status: StoreStatus::from((string) $doc['status']),
            notes: isset($doc['notes']) ? (string) $doc['notes'] : null,
            createdAt: Dates::fromBson($doc['created_at']),
            updatedAt: Dates::fromBson($doc['updated_at']),
        );
    }

    public function toDocument(): array
    {
        return [
            'store_uuid' => $this->storeUuid,
            'name'       => $this->name,
            'address'    => $this->address,
            'email'      => $this->email,
            'phone'      => $this->phone,
            'status'     => $this->status->value,
            'notes'      => $this->notes,
            'created_at' => new UTCDateTime($this->createdAt),
            'updated_at' => new UTCDateTime($this->updatedAt),
        ];
    }

    public function toArray(): array
    {
        return [
            'store_uuid' => $this->storeUuid,
            'name'       => $this->name,
            'address'    => $this->address,
            'email'      => $this->email,
            'phone'      => $this->phone,
            'status'     => $this->status->value,
            'notes'      => $this->notes,
            'created_at' => Dates::toIso($this->createdAt),
            'updated_at' => Dates::toIso($this->updatedAt),
        ];
    }
}
