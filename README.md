# MongoTest Simple Client

A simple Go application to demonstrate connecting to MongoDB Atlas and querying the `sample_mflix` database.

## Prerequisites

- Go 1.16 or higher
- Access to the internet (to connect to the MongoDB Atlas cluster)

## Configuration

The application uses a `config.json` file in the project root to store sensitive information like MongoDB credentials. A sample `config.json` looks like this:

```json
{
  "mongo_user": "YOUR_USERNAME",
  "mongo_password": "YOUR_PASSWORD",
  "mongo_scheme": "mongodb+srv",
  "mongo_host": "cluster0.dbpelmw.mongodb.net",
  "mongo_uri": "/?appName=Cluster0"
}
```

## How to Run

To run any of the applications, simply execute them using `go run`:

```bash
go run cmd/simpleclient/main.go
go run cmd/orders/main.go -file YOUR_FILE.csv
go run cmd/sales/main.go -file YOUR_FILE.csv
```

## Project Structure

- `cmd/simpleclient/main.go`: The main entry point for the application.
- `GEMINI.md`: Development standards and mandates for this project.
- data files in directory: ./data
