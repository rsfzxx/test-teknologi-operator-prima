<?php

declare(strict_types=1);

namespace App\Validation;

use App\Domain\StoreStatus;
use App\Exception\AppException;
use App\Repository\StoreFilter;

final class StoreValidator
{
    public const MAX_NAME_LENGTH    = 150;
    public const MAX_ADDRESS_LENGTH = 255;
    public const MAX_EMAIL_LENGTH   = 254;
    public const MAX_NOTES_LENGTH   = 500;
    public const MAX_SEARCH_LENGTH  = 100;
    public const MAX_PAGE           = 1_000_000;
    public const DEFAULT_PER_PAGE   = 10;
    public const MAX_PER_PAGE       = 100;

    private const STORE_FIELDS  = ['name', 'address', 'email', 'phone', 'status', 'notes'];
    private const STATUS_FIELDS = ['status'];

    private const PHONE_PATTERN = '/^\+?[0-9]{8,15}$/';

    /**
     *
     * @return array{name: string, address: string, email: string, phone: string, status: StoreStatus, notes: ?string}
     */
    public function validateStore(array $body): array
    {
        $this->rejectUnknownFields($body, self::STORE_FIELDS);

        $errors = [];
        $data = [
            'name'    => $this->requiredString($body, 'name', self::MAX_NAME_LENGTH, $errors),
            'address' => $this->requiredString($body, 'address', self::MAX_ADDRESS_LENGTH, $errors),
            'email'   => $this->email($body, $errors),
            'phone'   => $this->phone($body, $errors),
            'status'  => $this->status($body, $errors),
            'notes'   => $this->notes($body, $errors),
        ];

        if ($errors !== []) {
            throw AppException::validation($errors);
        }

        return $data;
    }

    public function validateStatusChange(array $body): StoreStatus
    {
        $this->rejectUnknownFields($body, self::STATUS_FIELDS);

        $errors = [];
        $status = $this->status($body, $errors);

        if ($errors !== [] || $status === null) {
            throw AppException::validation($errors);
        }

        return $status;
    }

    public function validateListQuery(array $query): StoreFilter
    {
        $errors = [];

        $page = $this->intQuery($query, 'page', 1, $errors);
        if ($page !== null && ($page < 1 || $page > self::MAX_PAGE)) {
            $errors[] = AppException::field('page', sprintf('must be between 1 and %d', self::MAX_PAGE));
        }

        $perPage = $this->intQuery($query, 'per_page', self::DEFAULT_PER_PAGE, $errors);
        if ($perPage !== null && ($perPage < 1 || $perPage > self::MAX_PER_PAGE)) {
            $errors[] = AppException::field('per_page', sprintf('must be between 1 and %d', self::MAX_PER_PAGE));
        }

        $search = trim($this->stringQuery($query, 'search', $errors));
        if ($this->length($search) > self::MAX_SEARCH_LENGTH) {
            $errors[] = AppException::field('search', sprintf('must be at most %d characters', self::MAX_SEARCH_LENGTH));
        }

        $status = null;
        $rawStatus = trim($this->stringQuery($query, 'status', $errors));
        if ($rawStatus !== '') {
            $status = StoreStatus::fromInput($rawStatus);
            if ($status === null) {
                $errors[] = AppException::field('status', 'must be one of: ' . implode(', ', StoreStatus::values()));
            }
        }

        $sortBy = strtolower(trim($this->stringQuery($query, 'sort_by', $errors)));
        $usingDefaultSort = $sortBy === '';
        if ($usingDefaultSort) {
            $sortBy = 'created_at';
        }
        if (!in_array($sortBy, StoreFilter::SORT_FIELDS, true)) {
            $errors[] = AppException::field('sort_by', 'must be one of: ' . implode(', ', StoreFilter::SORT_FIELDS));
        }

        $sortDesc = $usingDefaultSort;
        $order = strtolower(trim($this->stringQuery($query, 'order', $errors)));
        switch ($order) {
            case '':
                break;
            case 'asc':
                $sortDesc = false;
                break;
            case 'desc':
                $sortDesc = true;
                break;
            default:
                $errors[] = AppException::field('order', 'must be one of: asc, desc');
        }

        if ($errors !== []) {
            throw AppException::badRequest('Invalid query parameters', $errors);
        }

        return new StoreFilter(
            page: (int) $page,
            perPage: (int) $perPage,
            search: $search,
            status: $status,
            sortBy: $sortBy,
            sortDesc: $sortDesc,
        );
    }

    private function rejectUnknownFields(array $body, array $allowed): void
    {
        $errors = [];
        foreach (array_keys($body) as $key) {
            if (!in_array($key, $allowed, true)) {
                $errors[] = AppException::field((string) $key, 'is not allowed');
            }
        }

        if ($errors !== []) {
            throw AppException::badRequest('Invalid request body', $errors);
        }
    }

    private function requiredString(array $body, string $field, int $maxLength, array &$errors): string
    {
        $value = $body[$field] ?? null;

        if ($value === null) {
            $errors[] = AppException::field($field, 'is required');
            return '';
        }
        if (!is_string($value)) {
            $errors[] = AppException::field($field, 'must be a string');
            return '';
        }

        $value = trim($value);
        if ($value === '') {
            $errors[] = AppException::field($field, 'is required');
        } elseif ($this->length($value) > $maxLength) {
            $errors[] = AppException::field($field, sprintf('must be at most %d characters', $maxLength));
        }

        return $value;
    }

    private function email(array $body, array &$errors): string
    {
        $before = count($errors);
        $email = strtolower($this->requiredString($body, 'email', self::MAX_EMAIL_LENGTH, $errors));

        if (count($errors) === $before && filter_var($email, FILTER_VALIDATE_EMAIL) === false) {
            $errors[] = AppException::field('email', 'must be a valid email address');
        }

        return $email;
    }

    private function phone(array $body, array &$errors): string
    {
        $before = count($errors);
        $phone = $this->requiredString($body, 'phone', 30, $errors);
        if (count($errors) !== $before) {
            return $phone;
        }

        $phone = (string) preg_replace('/[\s\-().]/', '', $phone);
        if (preg_match(self::PHONE_PATTERN, $phone) !== 1) {
            $errors[] = AppException::field('phone', 'must be 8-15 digits and may start with +');
        }

        return $phone;
    }

    private function status(array $body, array &$errors): ?StoreStatus
    {
        $before = count($errors);
        $raw = $this->requiredString($body, 'status', 20, $errors);
        if (count($errors) !== $before) {
            return null;
        }

        $status = StoreStatus::fromInput($raw);
        if ($status === null) {
            $errors[] = AppException::field('status', 'must be one of: ' . implode(', ', StoreStatus::values()));
        }

        return $status;
    }

    private function notes(array $body, array &$errors): ?string
    {
        $value = $body['notes'] ?? null;
        if ($value === null) {
            return null;
        }
        if (!is_string($value)) {
            $errors[] = AppException::field('notes', 'must be a string');
            return null;
        }

        $value = trim($value);
        if ($value === '') {
            return null;
        }
        if ($this->length($value) > self::MAX_NOTES_LENGTH) {
            $errors[] = AppException::field('notes', sprintf('must be at most %d characters', self::MAX_NOTES_LENGTH));
        }

        return $value;
    }

    private function intQuery(array $query, string $key, int $default, array &$errors): ?int
    {
        $raw = $this->stringQuery($query, $key, $errors);
        if ($raw === '') {
            return $default;
        }
        if (preg_match('/^-?\d+$/', $raw) !== 1) {
            $errors[] = AppException::field($key, 'must be a number');
            return null;
        }

        return (int) $raw;
    }

    private function stringQuery(array $query, string $key, array &$errors): string
    {
        $value = $query[$key] ?? '';
        if (!is_string($value)) {
            $errors[] = AppException::field($key, 'must be a single value');
            return '';
        }

        return $value;
    }

    private function length(string $value): int
    {
        return function_exists('mb_strlen')
            ? mb_strlen($value, 'UTF-8')
            : (int) preg_match_all('/./su', $value);
    }
}
