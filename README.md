# MongoTest Trade Management Suite

A suite of Go applications for managing and querying trade and sales data in MongoDB Atlas.

## Prerequisites

- Go 1.20 or higher
- Access to the internet (to connect to the MongoDB Atlas cluster)

## Configuration

All applications use a `config.json` file in the project root for MongoDB credentials. Use `config.template.json` as a starting point.

```json
{
  "mongo_user": "YOUR_USERNAME",
  "mongo_password": "YOUR_PASSWORD",
  "mongo_scheme": "mongodb+srv",
  "mongo_host": "cluster0.dbpelmw.mongodb.net",
  "mongo_uri": "/?appName=Cluster0"
}
```

## Applications

### 1. Simple Client
A basic connectivity test that fetches movies from a sample database.
```bash
go run cmd/simpleclient/main.go
```

### 2. Orders & Sales Importer
Import orders or sales data from CSV or Excel files. Use the `-date-base` flag to correctly parse relative dates (e.g., "Sun, Feb 1").

```bash
# Import orders
go run cmd/orders/main.go -file path/to/orders.csv -date-base 2026-02

# Import sales
go run cmd/sales/main.go -file path/to/sales.csv -date-base 2026-02
```

### 3. Query Tool
Perform flexible JSON-based queries on your collections.
```bash
# Query for sales with qty > 10
go run cmd/query/main.go -collection sales -query '{"qty": {"$gt": 10}}'

# Query for orders on a specific date
go run cmd/query/main.go -collection orders -query '{"order_date": "2026-02-01T00:00:00Z"}'
```

## Project Structure

- `cmd/`: Application entry points (`simpleclient`, `orders`, `sales`, `query`).
- `internal/config/`: Configuration loading logic.
- `internal/models/`: Shared data structures (`Trade` model).
- `internal/mongodb/`: Database connection and repository logic.
- `internal/parser/`: CSV and Excel parsing utilities.
- `data/`: Sample data files.
- `GEMINI.md`: Development standards and mandates for this project.
