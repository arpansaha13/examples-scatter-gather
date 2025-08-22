# Scatter by Distribution: Combined Controller and Aggregator

This example implements a minimal scatter-gather pattern where the controller acts as the aggregator.

## Architecture

- **Gateway**: Acts as the aggregator/controller that makes parallel calls to microservices.
- **Service A**: Microservice that provides dummy data.
- **Service B**: Microservice that provides dummy data.

## How It Works

1. **Scatter**: The gateway makes parallel HTTP calls to both services using goroutines.
2. **Gather**: The gateway collects responses from all services and aggregates them.
3. **Response**: Returns a combined JSON response to the client.

## Endpoints

### Gateway (Port 4000)
- `GET /` - Home endpoint returning status message.
- `GET /data` - Gathers data from service A and service B.

### Service A (Port 4010)
- `GET /` - Home endpoint returning service status.
- `GET /data` - Returns dummy JSON data.

### Service B (Port 4020)
- `GET /` - Home endpoint returning service status.
- `GET /data` - Returns dummy JSON data.

## Running the System

### Using Docker Compose

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build

# Or run in watch mode
docker-compose watch --prune
```

### Manual Go Run

```bash
# Terminal 1 - Start Service A
cd service_a
go run main.go

# Terminal 2 - Start Service B
cd service_b
go run main.go

# Terminal 3 - Start Gateway
cd gateway
go run main.go
```

## Testing

Once all services are running, you can test the endpoints:

```bash
# Test gateway
curl http://localhost:4000/
curl http://localhost:4000/data

# Test service A
curl http://localhost:4010/
curl http://localhost:4010/data

# Test service B
curl http://localhost:4020/
curl http://localhost:4020/data
```

## Stopping the System

```bash
# If using Docker Compose
docker-compose down

# If running manually, use Ctrl+C in each terminal
```
