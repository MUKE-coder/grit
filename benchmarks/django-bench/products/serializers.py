from rest_framework import serializers

from .models import Product


class ProductSerializer(serializers.ModelSerializer):
    """
    Emits the same JSON types the other frameworks do.

    DRF renders DecimalField as a quoted string by default and BigIntegerField
    as a number; Grit, Laravel, Express and Bun all emit price as a JSON number.
    Coercing here keeps the payloads comparable — otherwise the response bodies
    differ in shape and this stops being a like-for-like test.

    deleted_at is excluded because no other framework returns it.
    """

    price = serializers.FloatField()
    stock = serializers.IntegerField()
    version = serializers.IntegerField()

    class Meta:
        model = Product
        fields = [
            "id",
            "name",
            "sku",
            "description",
            "price",
            "stock",
            "active",
            "version",
            "created_at",
            "updated_at",
        ]
