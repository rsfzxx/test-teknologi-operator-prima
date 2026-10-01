<?php

declare(strict_types=1);

namespace App\Repository;

use App\Domain\StoreStatus;

final class StoreFilter
{
    public const SORT_FIELDS = ['name', 'email', 'status', 'created_at', 'updated_at'];

    public function __construct(
        public readonly int $page,
        public readonly int $perPage,
        public readonly string $search,
        public readonly ?StoreStatus $status,
        public readonly string $sortBy,
        public readonly bool $sortDesc,
    ) {
    }

    public function offset(): int
    {
        return ($this->page - 1) * $this->perPage;
    }
}
