# Bumper

## Overview
This Go application provides a simple version bumping service. It includes two primary functionalities: bumping the version of a given semantic version string and retrieving the current version of the application. The application leverages the `go-semver` package for semantic version parsing and operations, and the `gorilla/mux` package for routing HTTP requests.

## Features
- **Version Bumping:** POST requests to `/bump` will increment the specified version according to the specified mode (major, minor, or patch).
- **Get Current Version:** GET requests to `/version` will return the current version of the application.

## Requirements
- Go programming language
- External Go packages:
    - `github.com/adamwasila/go-semver`
    - `github.com/gorilla/mux`

## Installation
To run this application, you need to have Go installed on your system. After installing Go, you can set up the project by following these steps:
1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run `go get` to install the required dependencies.

## Usage
To start the server, run the following command in the terminal from the project directory:
```bash
go run main.go
```

### Running with Docker
To run the application using Docker, follow these steps:
1. Build the Docker image:
   ```bash
    make docker-build
   ```
2. Run the Docker container:
   ```bash
    make docker-run
   ```
### Makefile Commands
- `make test`: Run tests.
- `make build`: Build the application.
- `make clean`: Clean the build directory.
- `make run ARGS="<args>"`: Run the application with specified arguments.
- `make docker-build`: Build the Docker image.
- `make docker-run`: Run the Docker container.
- `make docker-push`: Push the Docker image to a registry.

### Dockerfile
The Dockerfile uses a multi-stage build process:
- The first stage uses `golang:1.19` to build the application.
- The second stage uses `alpine:3.14` to create a lightweight container.

### Endpoints
- **Bump Version**
    - **URL:** `/bump`
    - **Method:** `POST`
    - **Body:**
      ```json
      {
        "version": "current version",
        "currentVersion": "current version to bump",
        "class": "major/minor/patch"
      }
      ```
    - **Response:**
      ```json
      {
        "statusCode": "0 for success, error codes for failure",
        "newVersion": "bumped version string"
      }
      ```

- **Get Current Version**
    - **URL:** `/version`
    - **Method:** `GET`
    - **Response:**
      ```json
      {
        "version": "current application version"
      }
      ```

## Contributing
Contributions to this project are welcome. Please ensure that you follow the established code style and add unit tests for new functionalities. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)

---

**Note:** This README provides a basic guide to using and contributing to this Go application. For more detailed documentation, refer to the comments within the code.
