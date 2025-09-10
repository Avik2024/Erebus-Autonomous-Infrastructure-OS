# Migrations Folder

This folder contains SQL migration files for the Erebus backend database.

**Naming convention:**

- `<version>_<description>.up.sql` → migration to apply changes
- `<version>_<description>.down.sql` → migration to revert changes
- Version numbers must be sequential (`001`, `002`, …)
- Always create `up` and `down` files together

**Example usage:**

```bash
# Apply all migrations
docker-compose run --rm migrate \
  -path=/migrations \
  -database "postgres://erebus:erebus@postgres:5432/erebus_db?sslmode=disable" up

# Revert last migration
docker-compose run --rm migrate \
  -path=/migrations \
  -database "postgres://erebus:erebus@postgres:5432/erebus_db?sslmode=disable" down
