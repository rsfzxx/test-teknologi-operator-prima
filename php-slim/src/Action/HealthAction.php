<?php

declare(strict_types=1);

namespace App\Action;

use App\Http\JsonResponder;
use MongoDB\Database;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Throwable;

final class HealthAction
{
    public function __construct(
        private readonly Database $database,
        private readonly JsonResponder $responder,
    ) {
    }

    public function __invoke(Request $request, Response $response): Response
    {
        try {
            $this->database->command(['ping' => 1]);
        } catch (Throwable $e) {
            error_log('health check failed: ' . $e->getMessage());

            return $this->responder->json($response, 503, [
                'success' => false,
                'message' => 'Service unavailable',
                'data'    => ['status' => 'unhealthy', 'database' => 'down'],
            ]);
        }

        return $this->responder->success($response, 200, 'Service is healthy', [
            'status'   => 'healthy',
            'database' => 'up',
        ]);
    }
}
