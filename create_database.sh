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

# CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
# uncomment if using psql version < 13

# Write these values into a .env file, or uncomment these lines
# DB_HOST="localhost"
# DB_PORT="5432"
# DB_NAME="lang_api"
# DB_USER="postgres"
# DB_PASSWORD="password"

# SQL statements to create the database and table
SQL_CREATE_DB="CREATE DATABASE $DB_NAME;"
SQL_CREATE_TABLE_USERS="CREATE TABLE users (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  profile TEXT NOT NULL,
  time_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"
SQL_CREATE_TABLE_CONTENT="CREATE TABLE content (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  author_id UUID REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  body_content TEXT,
  time_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"
SQL_CREATE_TABLE_SAVED_CONTENT="CREATE TABLE saved_content (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  content_id UUID REFERENCES content(id)
);"

SQL_CREATE_TABLE_CHAT_ROOMS="CREATE TABLE chat_room (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"
SQL_CREATE_CHAT_ROOM_PARTICIPANT="CREATE TABLE chat_room_participant (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID REFERENCES users(id) NOT NULL,
  chat_room_id UUID REFERENCES chat_room(id) NOT NULL
);"
SQL_CREATE_MESSAGES="CREATE TABLE messages (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  chat_room_id UUID REFERENCES chat_room(id),
  user_id UUID REFERENCES users(id),
  content TEXT NOT NULL,
  sent_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);"
SQL_INSERT_DATA="INSERT INTO users (id, username, profile) VALUES
  ('d2792a62-86a4-4c49-a909-b1e762c683a3', 'johnny_mnemonic', 'Discreetly transports sensitive data for corporations in a storage device implanted in his brain at the cost of his childhood memories.'),
  ('fc1b7d29-6aeb-432b-9354-7e4c65f15d4e', 'bob_loblaw', 'The Bluth familys new attorney'),
  ('9f0b1b5f-9cc5-4d14-aa9c-82cbe87e8a95', 'twinkle_toes', 'Air Nomad born in 12 BG and the Avatar during the Hundred Year War');
"

# Create the database
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "$SQL_CREATE_DB"

# Create the tables
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_USERS"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_CONTENT"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_SAVED_CONTENT"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_TABLE_CHAT_ROOMS"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_CHAT_ROOM_PARTICIPANT"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_CREATE_MESSAGES"

# Create default user
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_INSERT_DATA"

# Create mock data with csv file
# CSV_FILE="MOCK_DATA.csv"
# psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\COPY content FROM '$CSV_FILE' delimiter ',' CSV HEADER;"
# in psql you can also run this command: 
#  \copy content from '/Users/willwedmedyk/Downloads/MOCK_DATA.csv' delimiter ',' CSV HEADER;