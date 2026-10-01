<?php

declare(strict_types=1);

namespace App\Http;

use App\Exception\AppException;
use Psr\Http\Message\ResponseInterface;
use Psr\Http\Message\ServerRequestInterface;
use Slim\Exception\HttpException;
use Slim\Exception\HttpMethodNotAllowedException;
use Slim\Exception\HttpNotFoundException;
use Throwable;

final class HttpErrorHandler
{
    public function __construct(private readonly JsonResponder $responder)
    {
    }

    public function __invoke(
        ServerRequestInterface $request,
        Throwable $exception,
        bool $displayErrorDetails,
        bool $logErrors,
        bool $logErrorDetails,
    ): ResponseInterface {
        $allow = null;

        if ($exception instanceof AppException) {
            $appException = $exception;
        } elseif ($exception instanceof HttpMethodNotAllowedException) {
            $appException = AppException::methodNotAllowed();
            $allow = implode(', ', $exception->getAllowedMethods());
        } elseif ($exception instanceof HttpNotFoundException) {
            $appException = AppException::notFound('Endpoint not found');
        } elseif ($exception instanceof HttpException) {
            $status = $exception->getCode();
            $appException = $status >= 400 && $status < 600
                ? new AppException($status, $exception->getMessage())
                : AppException::internal($exception);
        } else {
            $appException = AppException::internal($exception);
        }

        if ($appException->status >= 500 && $logErrors) {
            $cause = $appException->getPrevious() ?? $appException;
            error_log(sprintf(
                'request failed: request_id=%s method=%s path=%s error=%s: %s at %s:%d',
                (string) $request->getAttribute('request_id', '-'),
                $request->getMethod(),
                $request->getUri()->getPath(),
                $cause::class,
                $cause->getMessage(),
                $cause->getFile(),
                $cause->getLine(),
            ));
        }

        $response = $this->responder->error($appException, $displayErrorDetails);

        return $allow !== null ? $response->withHeader('Allow', $allow) : $response;
    }
}
