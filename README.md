# MongoTest Trade Management Suite

A suite of Go applications for managing and querying trade and sales data in MongoDB Atlas.

## Prerequisites

- Go 1.20 or higher
- Access to the internet (to connect to the MongoDB Atlas cluster)
- Docker & Minikube (for API services)

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

## CLI Applications

### 1. Simple Client
A basic connectivity test.
```bash
go run cmd/simpleclient/main.go
```

### 2. Orders & Sales Importer
Import data from CSV/Excel. Use `-date-base` for relative dates (e.g., "Sun, Feb 1").
```bash
go run cmd/sales/main.go -file path/to/sales.csv -date-base 2026-02
```

### 3. Query & Diff Tools
```bash
# Query
go run cmd/query/main.go -query '{"order_date": "2026-02-01"}' -v

# Diff
go run cmd/diff/main.go -f1 file1.csv -f2 file2.csv
```

### 4. Reports Tool
```bash
go run cmd/reports/main.go -type summary -month 2026-02
```

## API Services & Minikube Deployment

The project includes two API services: `orders-api` and `reports-api`.

### 1. Build & Deploy
Point your shell to Minikube's Docker environment and build the images:

```bash
# Set docker env
eval $(minikube docker-env)

# Build Orders API
docker build -t orders-api:latest --build-arg APP_NAME=orders-api .

# Build Reports API
docker build -t reports-api:latest --build-arg APP_NAME=reports-api .

# Deploy to Minikube
kubectl apply -f k8s-manifests.yaml
```

### 2. Accessing the APIs
Use `minikube service` to get the URLs:
```bash
minikube service orders-api --url
minikube service reports-api --url
```

### 3. Testing with CURL
Once the services are running in Minikube, you can use `curl` to fetch data:

```bash
# Get an order by ID (replace URL with minikube service output)
curl $(minikube service orders-api --url)/api/orders/01LNK3Q

# Search for orders
curl "$(minikube service orders-api --url)/api/orders/search?q=Dropship"

# Get a reports summary for a date range
curl "$(minikube service reports-api --url)/api/reports/summary?startDate=2026-03-01&endDate=2026-03-31"
```

### 4. Endpoints

**Orders API (Internal: 8080, Service: 8083)**:
- `GET /api/orders/{orderId}`
- `GET /api/orders/tracking/{trackingId}`
- `GET /api/orders/search?q={term}`

**Reports API (Internal: 8080, Service: 8084)**:
- `GET /api/reports/summary?startDate=2026-01-01&endDate=2026-03-25`

## Project Structure

- `cmd/`: CLI entry points and API main files.
- `internal/`: Shared logic (config, models, mongodb, parser).
- `Dockerfile`: Multi-stage build for containerizing APIs.
- `k8s-manifests.yaml`: K8s deployment and service definitions.
- `GEMINI.md`: Project standards.
