# User CRUD API in Go

This is a basic User CRUD (Create, Read, Update, Delete) API built using Go with the Gin framework. It includes user authentication and health check endpoints.

## API Endpoints

### General Routes

- `GET /healthcheck` - Check if the server is running
- `POST /login` - User authentication endpoint

### User Routes

- `GET /users` - Retrieve all users
- `POST /users` - Create a new user
- `GET /users/:id` - Get a specific user by ID
- `PUT /users/:id` - Update a specific user by ID
- `DELETE /users/:id` - Delete a specific user by ID

## Prerequisites

- [MySQL](https://dev.mysql.com/downloads/) - For the database.
- [Go](https://go.dev/doc/install) - The programming language used for this API.

## Setup Instructions

1. **MySQL Setup**

   - Ensure MySQL is installed on your system
   - Create the database:
     ```sql
     CREATE DATABASE user_crud;
     USE DATABASE user_crud;
     ```
   - Log in to the MySQL CLI:
     ```bash
     mysql -u <username> -p
     ```
   - Run the `scripts.sql` file to create tables and insert initial records:
     ```bash
     source scripts.sql
     ```

2. **Environment Configuration**

   - Create a `.env` file in the root directory
   - Add the following environment variables:
     ```
     DBUSER=<database user name>
     DBPASS=<database password>
     JWT_SECRET_KEY=<jwt token secret key>
     ```

3. **Dependencies**

   - Ensure you have Go modules initialized (`go.mod` is included in the repo)
   - Install required dependencies:
     ```bash
     go mod tidy
     ```
   - This will install packages like `gin-gonic/gin`, `go-sql-driver/mysql`, `golang-jwt/jwt`, and others.

4. **Install Go Air for Live Reloading**

   - Since Gin doesn't detect code changes by default, we'll use the `air` package
   - Install air:
     ```bash
     go install github.com/air-verse/air@latest
     ```
   - Ensure `$GOPATH/bin` is added to your system's PATH environment variable
   - Verify installation:
     ```bash
     air -v
     ```

5. **Air Configuration**

   - An `air.toml` configuration file is already provided in the repository
   - To customize or generate a new one:
     - Edit `air.toml` as needed, or
     - Run `air init` to create a new configuration
   - For more details, check the [official Air documentation](https://github.com/air-verse/air)

6. **Run the Server**
   - Start the server with live reloading:
     ```bash
     air
     ```
   - The server will start on `localhost:8080`, and you can begin interacting with the API endpoints

## Usage

Once the server is running, you can test the API endpoints using tools like Postman, cURL, or any HTTP client. The server will automatically reload on code changes thanks to Air.

## Contributing

Feel free to submit issues or pull requests if you'd like to contribute to this project!
