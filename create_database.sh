#!/bin/bash

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please make sure PostgreSQL is installed."
    exit 1
fi

# PostgreSQL database credentials

# Load environment variables from .env file
if [ -f .env ]; then
    source .env
fi

# DB_HOST="localhost"
# DB_PORT="5432"
# DB_NAME="lang_api"
# DB_USER="postgres"
# DB_PASSWORD="password"

# SQL statements to create the database and table
SQL_CREATE_DB="CREATE DATABASE $DB_NAME;"
SQL_CREATE_TABLE_USERS="CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    time_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"
SQL_CREATE_TABLE_CONTENT="CREATE TABLE content (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    body_content TEXT,
    time_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"

SQL_INSERT_DATA="INSERT INTO users (username) VALUES
  ('johnny_mnemonic')
"

# Create the database
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "$SQL_CREATE_DB"

# Create the tables
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_USERS"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_CONTENT"

# Create default user
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_INSERT_DATA"

# Create mock data with csv file
# CSV_FILE="MOCK_DATA.csv"
# psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\COPY content FROM '$CSV_FILE' delimiter ',' CSV HEADER;"
# in psql you can also run this command: 
#  \copy content from '/Users/willwedmedyk/Downloads/MOCK_DATA.csv' delimiter ',' CSV HEADER;