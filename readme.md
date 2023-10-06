# MyServer Go HTTP API

This is a simple Go HTTP server that handles requests to manage `MyInfo` entities. It supports basic CRUD operations (Create, Read, Update, Delete) using appropriate HTTP methods (POST, GET, PUT, DELETE).

## Table of Contents

- [MyServer Go HTTP API](#myserver-go-http-api)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Installation and Setup](#installation-and-setup)
  - [Usage](#usage)
  - [API Endpoints](#api-endpoints)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)

## Overview

The server is a basic HTTP API implemented in Go (Golang) that allows the creation, retrieval, updating, and deletion of `MyInfo` entities. The server runs on port `8080` by default and uses a simple in-memory data store to manage the entities.

## Installation and Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/srikanthbhandary-teach/my-server.git
   cd my-server
   ```

2. Ensure you have Go installed. If not, [download and install Go](https://golang.org/dl/).

3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

To start the server, run the following command:

```bash
go run main.go
```

The server will start and be accessible at `http://localhost:8080`.

## API Endpoints

The API supports the following endpoints and methods:

- `POST /?id=<ID>`: Create a new `MyInfo` entity with the specified ID.
- `GET /?id=<ID>`: Retrieve a `MyInfo` entity by its ID.
- `PUT /?id=<ID>`: Update an existing `MyInfo` entity with the specified ID.
- `DELETE /?id=<ID>`: Delete a `MyInfo` entity by its ID.

## Testing

The project includes unit tests to verify the functionality of the server and its methods. To run the tests, execute the following command:

```bash
go test
```

## Contributing

If you would like to contribute to this project, feel free to open an issue or create a pull request. Contributions are welcome!

## License

This project is licensed under the [MIT License](LICENSE).