<?php

declare(strict_types=1);

namespace App\Action\Store;

use App\Http\JsonResponder;
use App\Service\StoreService;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;

/**
 * Menampilkan hasil aggregation pipeline (jumlah store per status).
 */
final class StoreSummaryAction
{
    public function __construct(
        private readonly StoreService $service,
        private readonly JsonResponder $responder,
    ) {
    }

    public function __invoke(Request $request, Response $response): Response
    {
        return $this->responder->success($response, 200, 'Store summary retrieved successfully', $this->service->summary());
    }
}
