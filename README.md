# Go and PostgreSQL Integration Project

## Description

This project demonstrates how to integrate a Go application with a PostgreSQL database.

## Prerequisites

* Go installed
* PostgreSQL installed and running
* [Optional] A tool like `direnv` to automatically load environment variables.

## Setup

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/eRaBxEs/task-manager-api
    cd https://github.com/eRaBxEs/task-manager-api
    ```

2.  **Create a .env file:**

    * Create a file named `.env` in the project root.
    * Add your database configuration to the `.env` file. See the `.env` file example below.
    * **Important:** Do *not* commit the `.env` file to your repository.

3.  **.env file content**

    ```
    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=your_database_name
    DB_USER=your_username
    DB_PASSWORD=your_password
    #   Add other environment variables as needed
    ```

4.  **Set up the database:**

    * Create a PostgreSQL database with the name specified in your `.env` file.
    * Ensure that the user specified in the `.env` file has the necessary permissions to access the database.


5.  **Run database migrations:**
    * Run the `run.sh` script.  This will automatically apply any pending database migrations.
    
6.  **Run the application:**

    * Run the `run.sh` script
    ```bash
    ./run.sh
    ```

## Running the application

    * The application will start, and you can access it at `http://localhost:8080` (or the port you have configured)

## Makefile (Alternative)

* A makefile is provided as an alternative to the `run.sh` script
    ```makefile
    DB_HOST := localhost
    DB_PORT := 5432
    DB_NAME := your_database_name
    DB_USER := your_username
    DB_PASSWORD := your_password

    run:
    export DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD
    go run cmd/myapp/main.go
    ```
* To use the makefile, simply run `make run`

## Testing

* To run the tests, use the following command:

    ```bash
    go test -v ./internal/handlers
    ```

## Explanation of Files

* `cmd/myapp/main.go`: The entry point of the Go application.
* `internal/config/config.go`: Handles loading of configuration from environment variables.
* `internal/database/db.go`: Handles the database connection.
* `internal/handlers/handlers.go`: Implements the HTTP handlers for the API.
* `internal/models/models.go`: Defines the data structures (e.g., Task).
* `internal/routes/routes.go`: Defines the API routes.
* `run.sh`: A shell script to run the application and execute database migrations.
* `.env`: A file that contains environment variables.
* `README.md`: This file, providing instructions on how to set up and run the application.
* `/migrations`:  Directory containing the goose migration files.

## Dependencies

* [github.com/gorilla/mux](https://github.com/gorilla/mux)

## License
