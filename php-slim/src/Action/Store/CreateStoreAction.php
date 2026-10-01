<?php

declare(strict_types=1);

namespace App\Action\Store;

use App\Http\JsonResponder;
use App\Http\RequestReader;
use App\Service\StoreService;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;

final class CreateStoreAction
{
    public function __construct(
        private readonly StoreService $service,
        private readonly JsonResponder $responder,
    ) {
    }

    public function __invoke(Request $request, Response $response): Response
    {
        $store = $this->service->create(RequestReader::jsonObject($request));

        return $this->responder
            ->success($response, 201, 'Store created successfully', $store->toArray())
            ->withHeader('Location', '/api/v1/stores/' . $store->storeUuid);
    }
}
