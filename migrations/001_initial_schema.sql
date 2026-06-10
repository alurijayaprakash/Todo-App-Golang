-- Migration: 001_initial_schema.sql
-- Description: Create initial todos table schema
-- Version: 1.0
-- Created: 2026-06-09

-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on created_at for better query performance
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at DESC);

-- Create index on completed status for filtering
CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
