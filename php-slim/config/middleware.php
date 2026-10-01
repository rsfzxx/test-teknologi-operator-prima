<?php

declare(strict_types=1);

use App\Http\HttpErrorHandler;
use App\Middleware\RequestIdMiddleware;
use Slim\App;

return static function (App $app): void {
    $container = $app->getContainer();
    $debug = (bool) $container->get('settings')['app']['debug'];

    $app->addRoutingMiddleware();

    $errorMiddleware = $app->addErrorMiddleware($debug, true, true);
    $errorMiddleware->setDefaultErrorHandler($container->get(HttpErrorHandler::class));

    $app->add(RequestIdMiddleware::class);
};
