#!/bin/bash

# Path to the file containing one JSON object per line
FILE="quotes.json"
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
auth_token=$(curl -X POST -d $login_data | grep Plaintext | grep -oE "\"[A-Z0-9]*?\"")

# Loop through each line in the file
while read -r req_body
do
  # Post the JSON in the current line as the request body
  curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $auth_token" -d "$req_body" $CREATE_QUOTES_ENDPOINT

  # Sleep for a short duration to avoid overwhelming the server (adjust according to rate limit)
  sleep 0.1

done < "$FILE"