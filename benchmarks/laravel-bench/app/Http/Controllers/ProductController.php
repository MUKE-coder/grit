<?php

namespace App\Http\Controllers;

use App\Models\Product;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

/**
 * The Laravel half of the Grit-vs-Laravel benchmark.
 *
 * Every choice here exists to match what Grit's generated handler does, so the
 * load test compares frameworks rather than implementations:
 *
 *   - the same default page size (20) and the same cap (100)
 *   - the same searchable columns and the same sortable allow-list
 *   - the same {data, meta:{total,page,page_size,pages}} envelope
 *   - the same {error:{code,message}} error envelope
 *
 * Written as plain Eloquent rather than API Resources. Resources would add a
 * transformation layer Grit's handler has no equivalent of, which would show up
 * as a Laravel tax that is really a difference in what the two are doing.
 */
class ProductController extends Controller
{
    private const DEFAULT_PAGE_SIZE = 20;

    private const MAX_PAGE_SIZE = 100;

    private const SEARCHABLE = ['name', 'sku', 'description'];

    private const SORTABLE = ['id', 'created_at', 'name', 'sku', 'description', 'stock'];

    public function index(Request $request): JsonResponse
    {
        $page = max(1, (int) $request->query('page', 1));

        $pageSize = (int) $request->query('page_size', self::DEFAULT_PAGE_SIZE);
        if ($pageSize < 1 || $pageSize > self::MAX_PAGE_SIZE) {
            $pageSize = self::DEFAULT_PAGE_SIZE;
        }

        $query = Product::query();

        if ($search = $request->query('search')) {
            $query->where(function ($q) use ($search) {
                foreach (self::SEARCHABLE as $column) {
                    // LIKE with LOWER() rather than ILIKE: portable across
                    // Postgres and SQLite, and the same expression Grit emits.
                    $q->orWhereRaw('LOWER('.$column.') LIKE LOWER(?)', ['%'.$search.'%']);
                }
            });
        }

        $sortBy = $request->query('sort_by', 'created_at');
        if (! in_array($sortBy, self::SORTABLE, true)) {
            $sortBy = 'created_at';
        }
        $sortOrder = strtolower((string) $request->query('sort_order', 'desc')) === 'asc' ? 'asc' : 'desc';

        $total = $query->count();

        $items = $query->orderBy($sortBy, $sortOrder)
            ->offset(($page - 1) * $pageSize)
            ->limit($pageSize)
            ->get();

        return response()->json([
            'data' => $items,
            'meta' => [
                'total' => $total,
                'page' => $page,
                'page_size' => $pageSize,
                'pages' => $pageSize > 0 ? (int) ceil($total / $pageSize) : 0,
            ],
        ]);
    }

    public function show(string $id): JsonResponse
    {
        $product = Product::find($id);

        if (! $product) {
            return $this->notFound();
        }

        return response()->json(['data' => $product]);
    }

    public function store(Request $request): JsonResponse
    {
        $validator = validator($request->all(), [
            'name' => 'required|string|max:255',
            'sku' => 'required|string|max:255',
            'description' => 'nullable|string',
            'price' => 'nullable|numeric',
            'stock' => 'nullable|integer',
            'active' => 'nullable|boolean',
        ]);

        if ($validator->fails()) {
            return response()->json([
                'error' => [
                    'code' => 'VALIDATION_ERROR',
                    'message' => $validator->errors()->first(),
                ],
            ], 422);
        }

        $product = Product::create($validator->validated());

        return response()->json([
            'data' => $product,
            'message' => 'Product created successfully',
        ], 201);
    }

    public function update(Request $request, string $id): JsonResponse
    {
        $product = Product::find($id);

        if (! $product) {
            return $this->notFound();
        }

        $validator = validator($request->all(), [
            'name' => 'sometimes|string|max:255',
            'sku' => 'sometimes|string|max:255',
            'description' => 'nullable|string',
            'price' => 'nullable|numeric',
            'stock' => 'nullable|integer',
            'active' => 'nullable|boolean',
        ]);

        if ($validator->fails()) {
            return response()->json([
                'error' => [
                    'code' => 'VALIDATION_ERROR',
                    'message' => $validator->errors()->first(),
                ],
            ], 422);
        }

        $product->fill($validator->validated());
        // Grit bumps version in a BeforeUpdate hook so offline clients can
        // detect a server-side change. Match it, or Laravel is doing less work.
        $product->version = $product->version + 1;
        $product->save();

        return response()->json([
            'data' => $product,
            'message' => 'Product updated successfully',
        ]);
    }

    public function destroy(string $id): JsonResponse
    {
        $product = Product::find($id);

        if (! $product) {
            return $this->notFound();
        }

        $product->delete();

        return response()->json(['message' => 'Product deleted successfully']);
    }

    private function notFound(): JsonResponse
    {
        return response()->json([
            'error' => [
                'code' => 'NOT_FOUND',
                'message' => 'Product not found',
            ],
        ], 404);
    }
}
