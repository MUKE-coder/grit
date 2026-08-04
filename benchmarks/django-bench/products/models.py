from django.db import models


class Product(models.Model):
    """
    Mapped onto the canonical `products` table from seed/schema.sql, which every
    framework in this benchmark shares.

    managed = False on purpose: the table is created once from that SQL file, so
    Django cannot quietly disagree with GORM or Eloquent about integer widths,
    timestamp precision or which columns get indexed. A benchmark where one side
    has an index the other lacks measures the index.
    """

    id = models.CharField(max_length=36, primary_key=True)
    name = models.CharField(max_length=255, null=True)
    sku = models.CharField(max_length=255, null=True)
    description = models.TextField(null=True)
    price = models.DecimalField(max_digits=12, decimal_places=2, null=True)
    stock = models.BigIntegerField(null=True)
    active = models.BooleanField(null=True)
    version = models.BigIntegerField(default=1)
    created_at = models.DateTimeField(null=True)
    updated_at = models.DateTimeField(null=True)
    deleted_at = models.DateTimeField(null=True)

    class Meta:
        db_table = "products"
        managed = False
