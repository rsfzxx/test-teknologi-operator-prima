<?php

declare(strict_types=1);

namespace App\Middleware;

use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Psr\Http\Server\MiddlewareInterface;
use Psr\Http\Server\RequestHandlerInterface;

final class RequestIdMiddleware implements MiddlewareInterface
{
    public const HEADER = 'X-Request-ID';

    public function process(ServerRequestInterface $request, RequestHandlerInterface $handler): ResponseInterface
    {
        $id = $request->getHeaderLine(self::HEADER);
        if (!self::isValid($id)) {
            $id = bin2hex(random_bytes(16));
        }

        $response = $handler->handle($request->withAttribute('request_id', $id));

        return $response->withHeader(self::HEADER, $id);
    }

    private static function isValid(string $id): bool
    {
        return $id !== '' && strlen($id) <= 64 && preg_match('/^[A-Za-z0-9_-]+$/', $id) === 1;
    }
}
