<?php

declare(strict_types=1);

use App\Database\Migrator;
use App\Database\MongoConnection;
use App\Database\PipelineLoader;
use App\Repository\MongoStoreRepository;
use App\Repository\StoreRepositoryInterface;
use DI\ContainerBuilder;
use Dotenv\Dotenv;
use MongoDB\Client;
use MongoDB\Database;
use Psr\Container\ContainerInterface;
use Psr\Http\Message\ResponseFactoryInterface;
use Slim\Psr7\Factory\ResponseFactory;

use function DI\autowire;

$dotenv = Dotenv::createImmutable(dirname(__DIR__));
$dotenv->safeLoad();
$dotenv->required(['MONGODB_URI', 'MONGODB_DATABASE'])->notEmpty();

$settings = require __DIR__ . '/settings.php';

$builder = new ContainerBuilder();
$builder->addDefinitions([
    'settings' => $settings,

    ResponseFactoryInterface::class => autowire(ResponseFactory::class),

    Client::class => static fn (): Client => MongoConnection::createClient($settings['mongodb']['uri']),

    Database::class => static fn (ContainerInterface $c): Database =>
        $c->get(Client::class)->selectDatabase($settings['mongodb']['database']),

    StoreRepositoryInterface::class => autowire(MongoStoreRepository::class),

    Migrator::class => static fn (ContainerInterface $c): Migrator =>
        new Migrator($c->get(Database::class), $settings['paths']['migrations']),

    PipelineLoader::class => static fn (): PipelineLoader =>
        new PipelineLoader($settings['paths']['pipelines']),
]);

return $builder->build();
