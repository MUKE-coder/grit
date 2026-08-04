<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

/**
 * Mirrors the Grit Product model column for column, so the benchmark compares
 * two frameworks doing the same work rather than two different schemas.
 */
class Product extends Model
{
    use HasUuids, SoftDeletes;

    protected $keyType = 'string';

    public $incrementing = false;

    protected $fillable = [
        'name',
        'sku',
        'description',
        'price',
        'stock',
        'active',
    ];

    protected $casts = [
        'price' => 'float',
        'stock' => 'integer',
        'active' => 'boolean',
        'version' => 'integer',
    ];

    protected $hidden = ['deleted_at'];
}
