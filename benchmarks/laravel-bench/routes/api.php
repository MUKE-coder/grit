<?php

use App\Http\Controllers\ProductController;
use Illuminate\Support\Facades\Route;

/*
 * Public CRUD, no auth, mounted at /api/v1/products to line up exactly with
 * Grit's route. Unauthenticated on both sides on purpose: with a token in play
 * the test would partly be measuring JWT parsing rather than the request path.
 */
Route::prefix('v1')->group(function () {
    Route::get('/products', [ProductController::class, 'index']);
    Route::get('/products/{id}', [ProductController::class, 'show']);
    Route::post('/products', [ProductController::class, 'store']);
    Route::put('/products/{id}', [ProductController::class, 'update']);
    Route::delete('/products/{id}', [ProductController::class, 'destroy']);
});
