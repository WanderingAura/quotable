# Introduction

Quotable is a backend JSON API that allows users to search for as well as create, edit and delete their favourite quotes!

A list of features the API current supports:

- User registration and authentication using auth tokens
- User email verification using smtp
- Request limiting based on IP address
- Basic CRUD operations such as CRUD on single quotes
- Advanced CRUD operations, including partial updates, text-based quote search with pagination and sorting, searching quotes by user
- User permissions so that unverified users can create quotes but cannot like them.

# Setup Instructions

## Setting up the database
This API uses a PostgreSQL database as part of its backend. So to run it we have to install postgres and set up the relevant databases.

First make sure that [PostgreSQL](https://www.postgresql.org/download/) and the [database migrations CLI tool](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) are installed on your computer.

Then execute the following commands:

`sudo -u postgres psql` to enter into postgres (or on Windows do `psql -U postgres` and enter the admin password you set on the postgres installation). Then

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

For convenience set the environment variable `export QUOTABLE_DB_DSN=postgres://test_user:test_pass@localhost/quotable` inside ~/.bashrc. (You don't have to do this if you replace $QUOTABLE_DB_DSN in the following instructions with the actual DSN). Remember to `source ~/.bashrc` in the relevant terminal afterwards.
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

## Environment variables

By default the app retrieves sensitive information from certain environment variables. Make a .envrc file in the project directory and fill it with the relevant usernames and passwords:
```bash
# Database credentials
export QUOTABLE_DB_DSN="***"
export QUOTABLE_TEST_DB_DSN="***"

# Credentials for an activated user
# This user will be used to post the quotes collected by the web scraper
export QUOTABLE_ADMIN_EMAIL="***"
export QUOTABLE_ADMIN_PASSWORD="***"

# SMTP credentials
export QUOTABLE_SMTP_USERNAME=***
export QUOTABLE_SMTP_PASSWORD=***
```

# Building and running the API locally
Make sure at least Go 1.22.1 is installed. In git bash execute `go get` to install the required go packages.

Type `make run/api` in the project directory and the API should start running on the default port 4000. (Make sure make is installed).

Then you can start to curl requests to the API.

# Endpoints

## Contents

| Type | Method | URL                    | Action                                   |
| ---- | ------ | ---------------------- | ---------------------------------------- |
| Healthcheck | GET | v1/version            | Get the version and deployment type of the app |
| User account | POST   | v1/user/register         | Register the user                        |
| User account | POST   | v1/users/password        | Change password of user (WIP)                  |
| User account | POST   | v1/tokens/auth           | Create an auth token for the user        |
| Query quotes | GET    | v1/quotes                | Query the quotes using url query params  |
| Query quotes | GET    | v1/quotes/:quote_id            | Query quote by quote ID, also outputs likes and dislikes   |
| Query quotes | GET    | v1/users/:user_id/quotes | Query the quotes of user with id user_id |
| Create/update quote | POST    | v1/quotes                | Creates a new quote as the authenticated user |
| Create/update quote | PATCH  | v1/quotes/:quote_id            | Update the quote  (WIP)                       |
| Delete quote | DELETE | v1/quotes/:quote_id            | Delete the quote                         |
| Like quote | POST | v1/quotes/:quote_id/like            | Like the quote as the authenticated user |

# Examples

## Check the API is running

Request:
`curl localhost:4000/v1/version`

Response:
```json
{
        "info": {
                "status": "available",
                "system_info": {
                        "environment": "development",
                        "version": "v1.0.0-0-g6a4c278-dirty"
                }
        }
}
```

## Register a user

Request:
`curl -d '{"username":"alice", "email":"alice@gmail.com", "password":"password"}' localhost:4000/v1/user/register`

Response:
```json
{
        "user": {
                "id": 2,
                "created_at": "2024-10-05T11:16:10+08:00",
                "username": "alice",
                "email": "alice@gmail.com",
                "activated": false
        }
}
```

## Get an auth token

Request:
`curl -d '{"email":"alice@gmail.com", "password":"password"}' localhost:4000/v1/tokens/auth`

Response:
```json
{
        "auth_token": {
                "Plaintext": "ODTTRDMXSGADHHYMLW6343BTEA",
                "Hash": "dzAosGqj+Ip276EY3hH9TcKJIdwB5b5iAMdrsWjB1OA=",
                "UserID": 2,
                "Expiry": "2024-10-06T11:18:48.8397265+08:00",
                "Scope": "auth"
        }
}
```

## Search for quotes

## Search quotes posted by a specific user

## Webscraper

The webscraper tool can be used to scrape quotes off goodreads and send them as quote creation requests to the API. Then account used for this is defined by the environment variables `QUOTABLE_ADMIN_EMAIL` and `QUOTABLE_ADMIN_PASSWORD`.

1. Execute quote_scraper.py (this requires python with beautifulsoup installed). This should create a `quotes.txt` file that contains the list of scraped quotes.
2. Create an account using the 
3. Execute `add_quotes.sh`. This will make a curl request to the API for each quote in `quotes.txt`.