<?php

declare(strict_types=1);

namespace App\Repository;

use App\Domain\Store;
use App\Domain\StoreStatus;

interface StoreRepositoryInterface
{
    /**
     * @return array{0: list<Store>, 1: int} daftar store dan total data
     */
    public function findAll(StoreFilter $filter): array;

    public function findByUuid(string $uuid): ?Store;

    public function emailExists(string $email, ?string $exceptUuid = null): bool;

    /** @throws DuplicateKeyException */
    public function create(Store $store): void;

    /**
     * @param array{name: string, address: string, email: string, phone: string, status: StoreStatus, notes: ?string} $fields
     * @throws DuplicateKeyException
     */
    public function update(string $uuid, array $fields): ?Store;

    public function updateStatus(string $uuid, StoreStatus $status): ?Store;

    public function delete(string $uuid): bool;

    /**
     * Menjalankan aggregation pipeline pada collection stores.
     *
     * @return list<array>
     */
    public function aggregate(array $pipeline): array;
}
