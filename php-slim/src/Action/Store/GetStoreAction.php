<?php

declare(strict_types=1);

namespace App\Action\Store;

use App\Http\JsonResponder;
use App\Http\RequestReader;
use App\Service\StoreService;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;

final class GetStoreAction
{
    public function __construct(
        private readonly StoreService $service,
        private readonly JsonResponder $responder,
    ) {
    }

    public function __invoke(Request $request, Response $response, array $args): Response
    {
        $store = $this->service->get(RequestReader::uuidParam($args));

        return $this->responder->success($response, 200, 'Store retrieved successfully', $store->toArray());
    }
}
