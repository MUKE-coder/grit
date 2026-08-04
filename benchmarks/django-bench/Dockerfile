# Django behind gunicorn, which is how Django is actually deployed. Not
# `manage.py runserver`, which is single-threaded, reloads on file change, and
# warns you not to use it in production — benchmarking it would be a strawman.
#
# gevent workers rather than sync: the workload is IO-bound on Postgres, and
# sync workers would block a whole worker per in-flight query. This is Django
# configured the way its own deployment docs recommend for this shape of load.
FROM python:3.13-slim

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends libpq5 \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt gevent==24.11.1

COPY . .

EXPOSE 8080

# Workers sized to the container's 4 CPUs by gunicorn's own formula
# (2 x CPU + 1), each with a gevent pool so a worker keeps serving while a query
# is in flight.
# Access logging stays OFF. Writing a line per request is real I/O and every
# other framework here has it disabled — turning it on for Django alone would
# tax it for something none of the others are paying.
CMD ["gunicorn", "bench.wsgi:application", \
     "--bind", "0.0.0.0:8080", \
     "--workers", "9", \
     "--worker-class", "gevent", \
     "--worker-connections", "200", \
     "--backlog", "2048", \
     "--keep-alive", "5", \
     "--error-logfile", "-", \
     "--log-level", "error"]
