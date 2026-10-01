<?php

declare(strict_types=1);

namespace App\Service;

use App\Database\PipelineLoader;
use App\Domain\Store;
use App\Domain\StoreStatus;
use App\Exception\AppException;
use App\Repository\DuplicateKeyException;
use App\Repository\StoreRepositoryInterface;
use App\Support\Dates;
use App\Validation\StoreValidator;
use Ramsey\Uuid\Uuid;

final class StoreService
{
    private const MSG_NOT_FOUND       = 'Store not found';
    private const MSG_DUPLICATE_EMAIL = 'Store with the same email already exists';
    private const SUMMARY_PIPELINE    = 'store_summary_by_status';

    public function __construct(
        private readonly StoreRepositoryInterface $repository,
        private readonly StoreValidator $validator,
        private readonly PipelineLoader $pipelines,
    ) {
    }

    /**
     * @return array{stores: list<Store>, total: int, page: int, per_page: int}
     */
    public function list(array $query): array
    {
        $filter = $this->validator->validateListQuery($query);
        [$stores, $total] = $this->repository->findAll($filter);

        return [
            'stores'   => $stores,
            'total'    => $total,
            'page'     => $filter->page,
            'per_page' => $filter->perPage,
        ];
    }

    public function get(string $uuid): Store
    {
        return $this->repository->findByUuid($uuid)
            ?? throw AppException::notFound(self::MSG_NOT_FOUND);
    }

    public function create(array $body): Store
    {
        $data = $this->validator->validateStore($body);
        $this->assertEmailAvailable($data['email']);

        $now = Dates::now();
        $store = new Store(
            storeUuid: Uuid::uuid4()->toString(),
            name: $data['name'],
            address: $data['address'],
            email: $data['email'],
            phone: $data['phone'],
            status: $data['status'],
            notes: $data['notes'],
            createdAt: $now,
            updatedAt: $now,
        );

        try {
            $this->repository->create($store);
        } catch (DuplicateKeyException) {
            throw $this->duplicateEmail();
        }

        return $store;
    }

    public function update(string $uuid, array $body): Store
    {
        $data = $this->validator->validateStore($body);

        $this->get($uuid);
        $this->assertEmailAvailable($data['email'], $uuid);

        try {
            $updated = $this->repository->update($uuid, $data);
        } catch (DuplicateKeyException) {
            throw $this->duplicateEmail();
        }

        return $updated ?? throw AppException::notFound(self::MSG_NOT_FOUND);
    }

    public function updateStatus(string $uuid, array $body): Store
    {
        $status = $this->validator->validateStatusChange($body);

        return $this->repository->updateStatus($uuid, $status)
            ?? throw AppException::notFound(self::MSG_NOT_FOUND);
    }

    public function delete(string $uuid): void
    {
        if (!$this->repository->delete($uuid)) {
            throw AppException::notFound(self::MSG_NOT_FOUND);
        }
    }

    public function summary(): array
    {
        $pipeline = $this->pipelines->load(self::SUMMARY_PIPELINE);
        $result = $this->repository->aggregate($pipeline)[0] ?? [];

        $rows = [];
        foreach ($result['by_status'] ?? [] as $row) {
            $rows[$row['status']] = $row;
        }

        $byStatus = [];
        foreach (StoreStatus::cases() as $status) {
            $row = $rows[$status->value] ?? null;
            $byStatus[] = [
                'status'            => $status->value,
                'total'             => (int) ($row['total'] ?? 0),
                'percentage'        => (float) ($row['percentage'] ?? 0),
                'latest_created_at' => Dates::toIso($row['latest_created_at'] ?? null),
            ];
        }

        return [
            'total'     => (int) ($result['total'] ?? 0),
            'by_status' => $byStatus,
        ];
    }

    private function assertEmailAvailable(string $email, ?string $exceptUuid = null): void
    {
        if ($this->repository->emailExists($email, $exceptUuid)) {
            throw $this->duplicateEmail();
        }
    }

    private function duplicateEmail(): AppException
    {
        return AppException::conflict(self::MSG_DUPLICATE_EMAIL, [
            AppException::field('email', 'is already used by another store'),
        ]);
    }
}
