from django.urls import path

from products import views

urlpatterns = [
    path("api/health", views.health),
    path("api/v1/products", views.products_collection),
    path("api/v1/products/<str:pk>", views.products_item),
]
