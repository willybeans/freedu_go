#!/bin/bash

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please make sure PostgreSQL is installed."
    exit 1
fi

# PostgreSQL database credentials
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="lang_api"
DB_USER="postgres"
DB_PASSWORD="password"

# SQL statements to create the database and table
SQL_CREATE_DB="CREATE DATABASE $DB_NAME;"
SQL_CREATE_TABLE="CREATE TABLE users (
    id INT,
    age INT,
    username TEXT,
    password TEXT,
    time_created INT,
    first_name TEXT,
    last_name TEXT,
    email TEXT
);"
SQL_CREATE_SEQUENCE="CREATE SEQUENCE id_seq
  START WITH 1
  INCREMENT BY 1
  MINVALUE 1
  NO MAXVALUE
  CACHE 1;
"

# Create the database
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "$SQL_CREATE_DB"

# Create the table
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE"

# Create Sequence
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_SEQUENCE"

