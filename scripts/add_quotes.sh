#!/bin/bash
set -x #echo on

# path to the file containing quotes in json form
FILE="./scripts/data/quotes.txt"
CREATE_QUOTES_ENDPOINT=http://localhost:4000/v1/quotes
AUTHENTICATION_ENDPOINT=http://localhost:4000/v1/tokens/auth

# load the admin account username and password
source .envrc

login_data=$(cat <<EOF
{
  "email": "$QUOTABLE_ADMIN_EMAIL",
  "password": "$QUOTABLE_ADMIN_PASSWORD"
}
EOF
)
auth_token=$(curl -X POST -d "$login_data" $AUTHENTICATION_ENDPOINT | grep Plaintext | grep -oE "\"[A-Z0-9]*?\"" | tr -d '"')

while read -r req_body
do
  # post the JSON in the current line as the request body
  curl -H "Authorization: Bearer $auth_token" -d "$req_body" $CREATE_QUOTES_ENDPOINT

  # sleep for a short duration to avoid overwhelming the server (adjust according to rate limit)
  sleep 0.5

done < "$FILE"