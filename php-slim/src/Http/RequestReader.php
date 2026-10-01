<?php

declare(strict_types=1);

namespace App\Http;

use App\Exception\AppException;
use JsonException;
use Psr\Http\Message\ServerRequestInterface;
use Ramsey\Uuid\Uuid;

final class RequestReader
{
    private const MAX_BODY_BYTES = 1024 * 1024; // 1 MB

    public static function uuidParam(array $args, string $name = 'id'): string
    {
        $value = $args[$name] ?? '';

        if (!is_string($value) || !Uuid::isValid($value)) {
            throw AppException::badRequest('Invalid path parameter', [
                AppException::field($name, 'must be a valid UUID'),
            ]);
        }

        return strtolower($value);
    }

    public static function jsonObject(ServerRequestInterface $request): array
    {
        $raw = (string) $request->getBody();

        if (strlen($raw) > self::MAX_BODY_BYTES) {
            throw AppException::payloadTooLarge();
        }

        $trimmed = ltrim($raw);
        if ($trimmed === '') {
            throw AppException::badRequest('Request body is required');
        }

        try {
            $data = json_decode($raw, true, 64, JSON_THROW_ON_ERROR);
        } catch (JsonException) {
            throw AppException::badRequest('Request body contains malformed JSON');
        }

        if (!is_array($data) || $trimmed[0] !== '{') {
            throw AppException::badRequest('Request body must contain a single JSON object');
        }

        return $data;
    }
}
