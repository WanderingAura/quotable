#!/bin/bash

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
auth_token=$(curl -iX POST -d "$login_data" $AUTHENTICATION_ENDPOINT | grep Plaintext | grep -oE "\"[A-Z0-9]*?\"" | tr -d '"')

echo '"'"Authorization: Bearer $auth_token"'"'