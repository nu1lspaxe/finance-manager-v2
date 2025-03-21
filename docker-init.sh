#!/bin/bash
set -e

export PGPASSWORD=$POSTGRESQL_PASSWORD

until psql -U postgres -c '\l'; do
  echo "Waiting for PostgreSQL to start..."
  sleep 2
done

DATABASES="bank_system fm_user fm_record fm_record_bank"

for db in $DATABASES; do
  echo "Creating database: $db"
  psql -U postgres -c "CREATE DATABASE \"$db\";"

  SQL_FILE="/docker-entrypoint-initdb.d/${db}.sql"
  
  if [ -f "$SQL_FILE" ]; then
    echo "Applying $SQL_FILE to $db"
    psql -U postgres -d "$db" -f "$SQL_FILE"
  else
    echo "Warning: $SQL_FILE not found for $db, skipping..."
  fi
done

echo "Databases initialized successfully!"