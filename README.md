# auth-rest-api

## Description
This is a REST API for authentication. It uses JWT for authentication and Postgres for storing user data and revocation data.

## Pre-requisites
1. Docker-compose should be installed. Refer to https://docs.docker.com/compose/install/

## Installation
1. Clone the repository : `git clone https://github.com/kasasunil/auth-rest-api.git`
2. Change directory to the project directory : `cd auth-rest-api`
3. Run this command to start the application : `docker-compose up --build`
4. The application will be running on `http://localhost:8080`


## Use cases/Endpoints
1. Signing up a user
    - **URL**: `http://localhost:8080/signup`
    - **Description**: This endpoint will create a user in the database. (We are storing password in plain text for simulation purpose. In real world we should store password in encrypted form.)
    - **Curl**:
     ```
        curl --location 'http://localhost:8080/public/signup' \
          --header 'Content-Type: application/json' \
          --data '{
          "email":"<Email>",
          "password":"<password>"
          }'
     ```
2. Signing in a user
      - **URL**: `http://localhost:8080/public/signin`
      - **Description**: This endpoint will return a token which will be used for authorization of private endpoints.We store this token in database and use it for authorization.So for simulation purpose user should store this token for using it for the below private endpoints. (Ideally user should not store token, but for simulation purpose I am using token only (We might have used email as well)).
      - **Curl**:
      ```
       curl --location 'http://localhost:8080/public/signin' \
            --header 'Content-Type: application/json' \
            --data '{
            "email":"<Email>",
            "password":"<password>"
            }'
   ```
3. Authorization of token
   - **Description**: This endpoint just returns details of user. For the below private endpoint we are authorizing(verifying token) the token using the middleware.
   - **URL**: `http://localhost:8080/private/user`
   - **Note**: Don't remove the `Bearer` keyword from the token. It should be in the format `Bearer <Token>`
   - **Curl**:
   ```
    curl --location 'http://localhost:8080/private/user' \
    --header 'Authorization: Bearer <Token that will come as a response in signin endpoint: `http://localhost:8080/public/signin`>'
    ```
4. Revocation of token
    - **Description**: We can revoke the token using the below endpoint. This is mainly a admin action, who can revoke the token from backend.The revoked token will be stored in the database and will be checked for every request. If any private request is using revoked token we don't allow that request.
    - **URL**: `http://localhost:8080/public/revoke_token`
    - **Curl**:
    ```
   curl --location 'http://localhost:8080/public/revoke_token' \
    --header 'Content-Type: application/json' \
    --data '{
    "token": "<Token that will come as a response in signin endpoint: `http://localhost:8080/public/signin`>"
    }'
    ```
   
5. Refresh token
    - **Description**: We can refresh the token using the below endpoint. This is mainly a user action, who can refresh the token before it's expiration. The refreshed token will be updated with the existing token of user. Only non expired token can be refreshed. Even if a token has some expiry time left but if it is in revoked list, then it can't be refreshed.
    - **URL**: `http://localhost:8080/private/refresh_token`
    - **Note**: Don't remove the `Bearer` keyword from the token. It should be in the format `Bearer <Token>`
    - **Curl**:
    ```
   curl --location 'http://localhost:8080/private/refresh_token' \
    --header 'Authorization: Bearer <Token that will come as a response in signin endpoint: `http://localhost:8080/public/signin`>'
    ```