# shortly

This is a backend API built with Go. It is designed for scalability and simplicity.It serves as a URL shortener service with features like generating and retrieving shortened URLs.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/tin3ga/shortly.git
   cd shortly
   ```

2. Install dependencies:

   ```bash
   go mod tidy

   ```

3. Configure .env file:

   ```shell
   PORT=8088
   DATABASE_URL=postgresql://<user>:<pass>@<host>/<database>
   URL=http://localhost:8088/api/v1/

   ```

### Configuration

1. Add a GOOSE_DBSTRING in makefile
2. Run make to set up migrations, dev dependencies etc.

   ```bash
   make help

   ```

### Usage

1. Run the application

   ```bash
   go run main.go

   ```

2. Swagger Documentation
   Access API documentation at:

   ```bash
      http://localhost:8088/swagger/index.html

   ```

### Deployment

For production, ensure you:

- Use a secure .env file
- Build the binary:

  ```bash
  go build -o app main.go

  ```

- Deploy the binary or container to your server or cloud provider.

## License

This project is licensed under the [MIT license][1].

&copy; 2025 tin3ga

[1]: LICENSE
