<?php

declare(strict_types=1);

namespace App\Http;

use App\Exception\AppException;
use JsonException;
use Psr\Http\Message\ResponseFactoryInterface;
use Psr\Http\Message\ResponseInterface;

final class JsonResponder
{
    private const JSON_FLAGS = JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_THROW_ON_ERROR;

    public function __construct(private readonly ResponseFactoryInterface $responseFactory)
    {
    }

    public function json(ResponseInterface $response, int $status, array $payload): ResponseInterface
    {
        try {
            $body = json_encode($payload, self::JSON_FLAGS);
        } catch (JsonException $e) {
            error_log('Failed to encode JSON response: ' . $e->getMessage());
            $status = 500;
            $body = '{"success":false,"message":"Internal server error"}';
        }

        $response->getBody()->write($body);

        return $response
            ->withStatus($status)
            ->withHeader('Content-Type', 'application/json; charset=utf-8');
    }

    public function success(ResponseInterface $response, int $status, string $message, mixed $data): ResponseInterface
    {
        return $this->json($response, $status, [
            'success' => true,
            'message' => $message,
            'data'    => $data,
        ]);
    }

    public function paginated(ResponseInterface $response, string $message, array $data, int $page, int $perPage, int $total): ResponseInterface
    {
        return $this->json($response, 200, [
            'success' => true,
            'message' => $message,
            'data'    => $data,
            'meta'    => [
                'page'        => $page,
                'per_page'    => $perPage,
                'total'       => $total,
                'total_pages' => $perPage > 0 ? (int) ceil($total / $perPage) : 0,
            ],
        ]);
    }

    public function error(AppException $e, bool $withDebug = false): ResponseInterface
    {
        $payload = [
            'success' => false,
            'message' => $e->getMessage(),
        ];

        if ($e->errors !== []) {
            $payload['errors'] = $e->errors;
        }

        $previous = $e->getPrevious();
        if ($withDebug && $previous !== null) {
            $payload['debug'] = [
                'exception' => $previous::class,
                'message'   => $previous->getMessage(),
                'file'      => $previous->getFile() . ':' . $previous->getLine(),
            ];
        }

        return $this->json($this->responseFactory->createResponse(), $e->status, $payload);
    }
}
