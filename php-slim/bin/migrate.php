<?php

declare(strict_types=1);

use App\Database\Migrator;

require __DIR__ . '/../vendor/autoload.php';

try {
    $container = require __DIR__ . '/../config/container.php';

    /** @var Migrator $migrator */
    $migrator = $container->get(Migrator::class);

    $applied = $migrator->run(static function (string $message): void {
        fwrite(STDOUT, $message . PHP_EOL);
    });

    fwrite(STDOUT, $applied === []
        ? 'Tidak ada migrasi baru.' . PHP_EOL
        : sprintf('Selesai, %d migrasi dijalankan.', count($applied)) . PHP_EOL);
    exit(0);
} catch (Throwable $e) {
    fwrite(STDERR, 'Migrasi gagal: ' . $e->getMessage() . PHP_EOL);
    exit(1);
}
