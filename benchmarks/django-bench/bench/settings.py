"""
Django settings for the benchmark.

Production-shaped: DEBUG off, no unnecessary middleware, connection pooling on.
The middleware list is trimmed to what a JSON API actually needs — sessions,
messages, CSRF and auth all cost time per request and none of them do anything
for an unauthenticated API. Leaving them in would measure Django's admin stack
rather than Django's request path.
"""

import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = os.environ.get("SECRET_KEY", "benchmark-only-not-a-real-secret")
DEBUG = False
ALLOWED_HOSTS = ["*"]

INSTALLED_APPS = [
    "django.contrib.contenttypes",
    "django.contrib.staticfiles",
    "rest_framework",
    "products",
]

# Only what a JSON API needs. CommonMiddleware handles APPEND_SLASH and
# Content-Length; everything else here would be dead weight per request.
MIDDLEWARE = [
    "django.middleware.common.CommonMiddleware",
]

ROOT_URLCONF = "bench.urls"
WSGI_APPLICATION = "bench.wsgi.application"
TEMPLATES = []

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.postgresql",
        "NAME": os.environ.get("DB_NAME", "bench_django"),
        "USER": os.environ.get("DB_USER", "bench"),
        "PASSWORD": os.environ.get("DB_PASSWORD", "bench"),
        "HOST": os.environ.get("DB_HOST", "postgres"),
        "PORT": os.environ.get("DB_PORT", "5432"),
        # Django 5's built-in pooling. Without it every request opens a new
        # connection and Postgres forks a backend — exactly the connection churn
        # this benchmark found in Grit, and it would cost Django just as much.
        "OPTIONS": {
            "pool": {
                "min_size": 4,
                "max_size": int(os.environ.get("DB_MAX_OPEN_CONNS", "25")),
            }
        },
    }
}

REST_FRAMEWORK = {
    # No browsable API renderer. It is a development convenience that costs
    # content negotiation on every request and nobody serves it in production.
    "DEFAULT_RENDERER_CLASSES": ["rest_framework.renderers.JSONRenderer"],
    "DEFAULT_PARSER_CLASSES": ["rest_framework.parsers.JSONParser"],
    "DEFAULT_AUTHENTICATION_CLASSES": [],
    "DEFAULT_PERMISSION_CLASSES": [],
    "UNAUTHENTICATED_USER": None,
}

USE_TZ = True
TIME_ZONE = "UTC"
DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"
STATIC_URL = "static/"

LOGGING = {
    "version": 1,
    "disable_existing_loggers": False,
    # Errors only. A log line per request is real I/O and would tax Django
    # against frameworks whose logging is off.
    "handlers": {"console": {"class": "logging.StreamHandler"}},
    "root": {"handlers": ["console"], "level": "ERROR"},
}
