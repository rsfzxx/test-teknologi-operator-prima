<?php

declare(strict_types=1);

namespace App\Action\Store;

use App\Domain\Store;
use App\Http\JsonResponder;
use App\Service\StoreService;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;

final class ListStoresAction
{
    public function __construct(
        private readonly StoreService $service,
        private readonly JsonResponder $responder,
    ) {
    }

    public function __invoke(Request $request, Response $response): Response
    {
        $result = $this->service->list($request->getQueryParams());

        return $this->responder->paginated(
            $response,
            'Stores retrieved successfully',
            array_map(static fn (Store $s): array => $s->toArray(), $result['stores']),
            $result['page'],
            $result['per_page'],
            $result['total'],
        );
    }
}
