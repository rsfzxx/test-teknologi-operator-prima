<?php

declare(strict_types=1);

use App\Action\HealthAction;
use App\Action\Store\CreateStoreAction;
use App\Action\Store\DeleteStoreAction;
use App\Action\Store\GetStoreAction;
use App\Action\Store\ListStoresAction;
use App\Action\Store\StoreSummaryAction;
use App\Action\Store\UpdateStoreAction;
use App\Action\Store\UpdateStoreStatusAction;
use Slim\App;
use Slim\Routing\RouteCollectorProxy;

return static function (App $app): void {
    $app->get('/health', HealthAction::class);

    $app->group('/api/v1', function (RouteCollectorProxy $group): void {
        $group->get('/stores', ListStoresAction::class);
        $group->post('/stores', CreateStoreAction::class);

        $group->get('/stores/summary', StoreSummaryAction::class);

        $group->get('/stores/{id}', GetStoreAction::class);
        $group->put('/stores/{id}', UpdateStoreAction::class);
        $group->patch('/stores/{id}/status', UpdateStoreStatusAction::class);
        $group->delete('/stores/{id}', DeleteStoreAction::class);
    });
};