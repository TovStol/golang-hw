#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

echo "PostgreSQL is up - executing migrations"

PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<EOF
CREATE TABLE IF NOT EXISTS event (
    id            bigint primary key GENERATED ALWAYS AS IDENTITY,
    title         text        not null,
    date_time     timestamptz not null,
    end_date_time timestamptz not null,
    description   text,
    user_id       bigint      not null,
    notify_before bigint
);

CREATE TABLE IF NOT EXISTS notification (
    event_id   bigint      primary key,
    title      text        not null,
    event_date timestamptz not null,
    user_id    bigint      not null
);
EOF

echo "Migrations completed successfully"
