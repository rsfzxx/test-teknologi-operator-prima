<?php

declare(strict_types=1);

namespace App\Repository;

use App\Domain\Store;
use App\Domain\StoreStatus;
use App\Support\Dates;
use MongoDB\BSON\Regex;
use MongoDB\BSON\UTCDateTime;
use MongoDB\Collection;
use MongoDB\Database;
use MongoDB\Driver\Exception\Exception as DriverException;
use MongoDB\Driver\Exception\WriteException;
use MongoDB\Operation\FindOneAndUpdate;

final class MongoStoreRepository implements StoreRepositoryInterface
{
    public const COLLECTION = 'stores';

    private const DUPLICATE_KEY_CODE = 11000;

    private const PROJECTION = ['_id' => 0];

    private readonly Collection $collection;

    public function __construct(Database $database)
    {
        $this->collection = $database->selectCollection(self::COLLECTION);
    }

    public function findAll(StoreFilter $filter): array
    {
        $query = $this->buildQuery($filter);

        $total = $this->collection->countDocuments($query);
        if ($total === 0) {
            return [[], 0];
        }

        $cursor = $this->collection->find($query, [
            'projection' => self::PROJECTION,
            'sort'       => [$filter->sortBy => $filter->sortDesc ? -1 : 1, 'store_uuid' => 1],
            'skip'       => $filter->offset(),
            'limit'      => $filter->perPage,
        ]);

        $stores = [];
        foreach ($cursor as $doc) {
            $stores[] = Store::fromDocument($doc);
        }

        return [$stores, $total];
    }

    public function findByUuid(string $uuid): ?Store
    {
        $doc = $this->collection->findOne(['store_uuid' => $uuid], ['projection' => self::PROJECTION]);

        return $doc === null ? null : Store::fromDocument($doc);
    }

    public function emailExists(string $email, ?string $exceptUuid = null): bool
    {
        $query = ['email' => $email];
        if ($exceptUuid !== null) {
            $query['store_uuid'] = ['$ne' => $exceptUuid];
        }

        return $this->collection->countDocuments($query, ['limit' => 1]) > 0;
    }

    public function create(Store $store): void
    {
        try {
            $this->collection->insertOne($store->toDocument());
        } catch (DriverException $e) {
            throw $this->translate($e);
        }
    }

    public function update(string $uuid, array $fields): ?Store
    {
        return $this->findAndSet($uuid, [
            'name'    => $fields['name'],
            'address' => $fields['address'],
            'email'   => $fields['email'],
            'phone'   => $fields['phone'],
            'status'  => $fields['status']->value,
            'notes'   => $fields['notes'],
        ]);
    }

    public function updateStatus(string $uuid, StoreStatus $status): ?Store
    {
        return $this->findAndSet($uuid, ['status' => $status->value]);
    }

    public function delete(string $uuid): bool
    {
        return $this->collection->deleteOne(['store_uuid' => $uuid])->getDeletedCount() === 1;
    }

    public function aggregate(array $pipeline): array
    {
        return iterator_to_array($this->collection->aggregate($pipeline), false);
    }

    private function findAndSet(string $uuid, array $set): ?Store
    {
        $set['updated_at'] = new UTCDateTime(Dates::now());

        try {
            $doc = $this->collection->findOneAndUpdate(
                ['store_uuid' => $uuid],
                ['$set' => $set],
                [
                    'projection'     => self::PROJECTION,
                    'returnDocument' => FindOneAndUpdate::RETURN_DOCUMENT_AFTER,
                ],
            );
        } catch (DriverException $e) {
            throw $this->translate($e);
        }

        return $doc === null ? null : Store::fromDocument($doc);
    }

    private function buildQuery(StoreFilter $filter): array
    {
        $query = [];

        if ($filter->search !== '') {
            $regex = new Regex(preg_quote($filter->search), 'i');
            $query['$or'] = [
                ['name' => $regex],
                ['email' => $regex],
                ['phone' => $regex],
                ['address' => $regex],
            ];
        }

        if ($filter->status !== null) {
            $query['status'] = $filter->status->value;
        }

        return $query;
    }

    private function translate(DriverException $e): \Throwable
    {
        if ($e instanceof WriteException) {
            foreach ($e->getWriteResult()->getWriteErrors() as $writeError) {
                if ($writeError->getCode() === self::DUPLICATE_KEY_CODE) {
                    return new DuplicateKeyException($writeError->getMessage(), 0, $e);
                }
            }
        }

        if ($e->getCode() === self::DUPLICATE_KEY_CODE) {
            return new DuplicateKeyException($e->getMessage(), 0, $e);
        }

        return $e;
    }
}
