# Introduction

Quotable is a backend JSON API that allows users to search for as well as create, edit and delete their favourite quotes!

A list of features the API current supports:

- User registration and authentication using auth tokens
- User email verification using smtp
- Request limiting based on IP address
- Basic CRUD operations such as CRUD on single quotes
- Advanced CRUD operations, including partial updates, text-based quote search with pagination and sorting, searching quotes by user
- User permissions so that unverified users cannot create quotes but only read them

# Setup Instructions

## Setting up the database
This API uses a PostgreSQL database as part of its backend. So to run it we have to set it up.

First make sure that [PostgreSQL](https://www.postgresql.org/download/) and the [database migrations CLI tool](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) are installed on your computer.

Then execute the following commands:

`sudo -u postgres psql` to enter into postgres. Then

```
CREATE DATABASE quotable;
\c quotable
CREATE ROLE test_user WITH LOGIN PASSWORD 'test_pass';
CREATE EXTENSION IF NOT EXISTS citext;
```
If on postgres 15 or later also execute:
```
GRANT ALL PRIVILEGES ON DATABASE quotable TO test_user;
GRANT ALL ON SCHEMA public TO quotable;
```
This sets up the database and the user that we are going to use to connect to it.

## Setting up the tables

For convenience set the environment variable `export QUOTABLE_DSN=postgres://test_user:test_pass@localhost/quotable` inside ~/.bashrc. (You don't have to do this, but if you don't replace $QUOTABLE_DSN in the following instructions with the actual DSN). Remember to `source ~/.bashrc` in the relevant terminal afterwards.
Again make sure the migrate CLI tool is installed then execute the following in the project directory:
`migrate -path=./migrations -database=$QUOTABLE_DSN up`
If successful the output should be similar to the following:
```
1/u create_users_table (41.837151ms)
2/u add_users_indexes (64.278511ms)
3/u create_quotes_table (94.12246ms)
4/u add_quotes_check_constraints (113.006448ms)
5/u add_quotes_indexes (130.975954ms)
6/u add_quotes_triggers (149.255275ms)
7/u create_tokens_table (178.395178ms)
8/u add_permissions (214.201982ms)
```

The database is now fully set up.

## Running the API locally

Type `go run ./cmd/api -db-dsn=$QUOTABLE_DSN` in the project directory and the API should start running on the default port 4000,

Then you can start to curl requests to the API.
