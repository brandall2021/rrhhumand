#!/bin/sh
set -e

echo "=== RRHHumand Entrypoint ==="

if [ -n "$DATABASE_HOST" ]; then
  echo "Waiting for PostgreSQL..."
  until pg_isready -h "$DATABASE_HOST" -p "${DATABASE_PORT:-5432}" -U "${DATABASE_USER:-postgres}" 2>/dev/null; do
    sleep 2
  done
  echo "PostgreSQL is ready."

  echo "Running migrations..."
  for dir in /app/migrations/*/; do
    up_file="${dir}up.sql"
    if [ -f "$up_file" ]; then
      echo "  Applying: $up_file"
      PGPASSWORD="${DATABASE_PASSWORD}" psql \
        -h "$DATABASE_HOST" \
        -p "${DATABASE_PORT:-5432}" \
        -U "${DATABASE_USER:-postgres}" \
        -d "${DATABASE_NAME:-rrhhumand}" \
        -f "$up_file" 2>/dev/null || echo "  WARNING: migration may have failed (might be already applied)"
    fi
  done
  echo "Migrations complete."
fi

echo "Starting server..."
exec "$@"
