# Company API

Company API - demo of a REST API for managing company data.


# Setting up

1. Clone the repository

```bash
git clone https://github.com/NikMohorko/company-api
```

2. The repository contains the configuration file `.env.example`, which you should rename to `.env`. The file already contains example values that can be used for testing or replaced with new ones.

3. Navigate to project folder and build the Docker container:
```bash
docker build .
```

4. Start the container:
```bash
docker compose up
```
Once database initialization is finished, the container will be listening for requests on port provided in the configuration file (localhost:8080 by default).

# Use

Complete API documentation can be found [here](https://app.swaggerhub.com/apis-docs/NikMohorko/companyAPI/1.0#/).

1. Create a new user by calling the /user/create endpoint.
2. Authenticate through the /user/authenticate endpoint. You will receive a token that you can use in all other requests using Bearer authentication.


# Testing

When the container is running you can run integration tests with:
```bash
go test -v
```
After running tests you can remove data that was inserted into database with:
```bash
docker-compose down --rmi all --volumes
```

# Built With
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
