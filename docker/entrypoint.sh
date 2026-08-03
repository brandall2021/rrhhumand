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
      output=$(PGPASSWORD="${DATABASE_PASSWORD}" psql \
        -h "$DATABASE_HOST" \
        -p "${DATABASE_PORT:-5432}" \
        -U "${DATABASE_USER:-postgres}" \
        -d "${DATABASE_NAME:-rrhhumand}" \
        -v ON_ERROR_STOP=0 \
        -f "$up_file" 2>&1 || true)
      if printf '%s' "$output" | grep -q "^psql:.*ERROR:"; then
        echo "  !!! Migration reported errors (check below):"
        printf '%s\n' "$output" | grep "^psql:.*ERROR:" | head -10
      fi
    fi
  done
  echo "Migrations complete."
fi

echo "Starting server..."
exec "$@"
