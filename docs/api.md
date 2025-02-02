# API

## Locally

The API documentation is build with Swagger and is available
under http://localhost:8080/swagger/index.html.

If you have changed the default server port ``8080``, make sure to adjust it in the above URL.

## Endpoints

Currently, the following endpoints are available. More might come:

- /accounts
- /consume
- /coffee

### /accounts

Endpoint for creating accounts available in Coffy Machine. Consume requests need a referenced
account ID to charge an account. 

### /consume

Endpoint for charging accounts with coffee consumption.

### /coffee

Endpoint for creating or updating available coffee products in Coffy Machine.
