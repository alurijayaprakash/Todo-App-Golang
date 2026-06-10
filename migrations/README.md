# Database Migrations

This directory contains version-based SQL migration scripts for the Todo App database schema.

## Migration Naming Convention

Migrations follow the pattern: `XXX_description.sql`
- `XXX`: Three-digit version number (001, 002, 003, etc.)
- `description`: Brief description of the migration

## How to Apply Migrations

### Manual Execution (Recommended for Production)

Connect to your PostgreSQL database and execute the SQL files in order:

```bash
# Using psql CLI
psql -h <host> -U <user> -d <database> -f migrations/001_initial_schema.sql

# Or using Docker
docker exec -i <container_name> psql -U <user> -d <database> < migrations/001_initial_schema.sql
```

### Connection String Format

The application constructs the connection string automatically from these environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=todo_db
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSLMODE=disable
```

Or via command-line flags:
```bash
go run cmd/server/main.go \
  -db-host localhost \
  -db-port 5432 \
  -db-name todo_db \
  -db-user postgres \
  -db-password securepassword \
  -db-sslmode disable
```

## Migration Files

- **001_initial_schema.sql**: Creates the initial `todos` table with indexes

## Notes

- Each migration is idempotent (safe to run multiple times)
- Migrations include `IF NOT EXISTS` clauses to prevent errors on rerun
- Indexes are included for optimal query performance
- Always test migrations in a staging environment before production

## Rollback Strategy

For rollback, you would need corresponding `.down.sql` files. Currently, this project assumes forward-only migrations.
