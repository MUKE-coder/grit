"""
The Django side of the benchmark, written with DRF's function-based api_view.

Deliberately not a ModelViewSet. A viewset would add routing, permission and
content-negotiation machinery that Grit's generated handler has no equivalent
of, and that cost would read as a Django tax that is really a difference in what
the two are doing. This is the same four endpoints, doing the same work.

Every choice matches Grit's generated handler: same default page size (20) and
cap (100), same searchable columns, same sortable allow-list, the same
{data, meta} envelope, the same {error:{code,message}} shape, and the same
version bump on update.
"""

import math
import uuid

from django.db.models import Q
from django.db.models.functions import Lower
from django.utils import timezone
from rest_framework.decorators import api_view
from rest_framework.response import Response

from .models import Product
from .serializers import ProductSerializer

DEFAULT_PAGE_SIZE = 20
MAX_PAGE_SIZE = 100
SEARCHABLE = ["name", "sku", "description"]
SORTABLE = {"id", "created_at", "name", "sku", "description", "stock"}


def _int(value, fallback):
    try:
        return int(value)
    except (TypeError, ValueError):
        return fallback


def _error(code, message, status):
    return Response({"error": {"code": code, "message": message}}, status=status)


def _not_found():
    return _error("NOT_FOUND", "Product not found", 404)


@api_view(["GET"])
def health(_request):
    return Response({"status": "ok"})


@api_view(["GET", "POST"])
def products_collection(request):
    if request.method == "POST":
        return _create(request)
    return _list(request)


@api_view(["GET", "PUT", "DELETE"])
def products_item(request, pk):
    product = Product.objects.filter(pk=pk, deleted_at__isnull=True).first()
    if product is None:
        return _not_found()

    if request.method == "PUT":
        return _update(request, product)
    if request.method == "DELETE":
        return _destroy(product)
    return Response({"data": ProductSerializer(product).data})


def _list(request):
    page = max(1, _int(request.query_params.get("page"), 1))

    page_size = _int(request.query_params.get("page_size"), DEFAULT_PAGE_SIZE)
    if page_size < 1 or page_size > MAX_PAGE_SIZE:
        page_size = DEFAULT_PAGE_SIZE

    queryset = Product.objects.filter(deleted_at__isnull=True)

    search = request.query_params.get("search")
    if search:
        # LOWER(col) LIKE LOWER(%term%), matching what Grit emits. Not
        # icontains-per-column-ORed differently: the generated SQL has to be the
        # same shape or the databases do different work.
        condition = Q()
        for column in SEARCHABLE:
            condition |= Q(**{f"{column}__icontains": search})
        queryset = queryset.filter(condition)

    sort_by = request.query_params.get("sort_by")
    if sort_by not in SORTABLE:
        sort_by = "created_at"
    descending = (request.query_params.get("sort_order") or "desc").lower() != "asc"
    ordering = f"-{sort_by}" if descending else sort_by

    total = queryset.count()
    offset = (page - 1) * page_size
    rows = queryset.order_by(ordering)[offset : offset + page_size]

    return Response(
        {
            "data": ProductSerializer(rows, many=True).data,
            "meta": {
                "total": total,
                "page": page,
                "page_size": page_size,
                "pages": math.ceil(total / page_size) if page_size else 0,
            },
        }
    )


def _create(request):
    body = request.data or {}

    if not body.get("name"):
        return _error("VALIDATION_ERROR", "name is required", 422)
    if not body.get("sku"):
        return _error("VALIDATION_ERROR", "sku is required", 422)

    now = timezone.now()
    product = Product(
        id=str(uuid.uuid4()),
        name=body.get("name"),
        sku=body.get("sku"),
        description=body.get("description") or "",
        price=body.get("price") or 0,
        stock=body.get("stock") or 0,
        active=bool(body.get("active", False)),
        version=1,
        created_at=now,
        updated_at=now,
    )
    # force_insert because the primary key is set by hand; without it Django
    # issues an UPDATE first to find out whether the row exists, which is a
    # wasted round trip on every create.
    product.save(force_insert=True)

    return Response(
        {"data": ProductSerializer(product).data, "message": "Product created successfully"},
        status=201,
    )


def _update(request, product):
    body = request.data or {}

    for field in ("name", "sku", "description", "price", "stock", "active"):
        if field in body and body[field] is not None:
            setattr(product, field, body[field])

    product.version = (product.version or 0) + 1
    product.updated_at = timezone.now()
    product.save()

    return Response(
        {"data": ProductSerializer(product).data, "message": "Product updated successfully"}
    )


def _destroy(product):
    product.deleted_at = timezone.now()
    product.save(update_fields=["deleted_at"])
    return Response({"message": "Product deleted successfully"})
