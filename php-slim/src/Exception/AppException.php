<?php

declare(strict_types=1);

namespace App\Exception;

use RuntimeException;
use Throwable;

final class AppException extends RuntimeException
{
    /**
     * @param list<array{field: string, message: string}> $errors
     */
    public function __construct(
        public readonly int $status,
        string $message,
        public readonly array $errors = [],
        ?Throwable $previous = null,
    ) {
        parent::__construct($message, $status, $previous);
    }

    public static function badRequest(string $message, array $errors = []): self
    {
        return new self(400, $message, $errors);
    }

    public static function notFound(string $message): self
    {
        return new self(404, $message);
    }

    public static function methodNotAllowed(): self
    {
        return new self(405, 'Method not allowed');
    }

    public static function conflict(string $message, array $errors = []): self
    {
        return new self(409, $message, $errors);
    }

    public static function payloadTooLarge(): self
    {
        return new self(413, 'Request body too large');
    }

    public static function validation(array $errors): self
    {
        return new self(422, 'Validation failed', $errors);
    }

    public static function internal(Throwable $previous): self
    {
        return new self(500, 'Internal server error', [], $previous);
    }

    public static function field(string $field, string $message): array
    {
        return ['field' => $field, 'message' => $message];
    }
}
