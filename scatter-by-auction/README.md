# Scatter by Auction: RabbitMQ-Based Scatter-Gather Pattern

This example implements a minimal scatter-gather pattern using RabbitMQ as the message broker, where the gateway acts as the controller and aggregator.

## Architecture

- **Gateway**: Acts as the controller/aggregator that publishes messages to RabbitMQ and consumes responses.
- **RabbitMQ**: Message broker with a fanout exchange for broadcasting messages to all services.
- **Service A**: Microservice that consumes messages and publishes responses to the results queue.
- **Service B**: Microservice that consumes messages and publishes responses to the results queue.

## How It Works

1. **Scatter**: The gateway publishes a message with a unique correlation ID to a RabbitMQ fanout exchange.
2. **Broadcast**: Both services receive the message from the exchange and process it independently.
3. **Gather**: Services publish their responses to a "results" queue with the correlation ID.
4. **Aggregation**: The gateway consumes from the results queue and routes responses to the appropriate buffered channels.
5. **Response**: The gateway waits for responses (with timeout) and returns aggregated results to the client.

## Components

### Gateway (Port 4000)
- **Global sync.Map**: Stores correlation IDs mapped to buffered channels for each request.
- **HTTP Server**: Exposes endpoints for home and data requests.
- **RabbitMQ Publisher**: Publishes messages to the scatter-gather exchange.
- **Aggregator**: Runs in a separate goroutine, consuming from the results queue.
- **Response Handler**: Waits for responses with timeout and minimum count requirements.

### RabbitMQ
- **Exchange**: "scatter-gather" (fanout type) for broadcasting messages.
- **Service Queues**: Individual queues for each service bound to the exchange.
- **Results Queue**: Central queue where all services publish their responses.

### Service A (Port 4010) & Service B (Port 4020)
- **HTTP Server**: Exposes home and data endpoints.
- **RabbitMQ Consumer**: Consumes messages from their respective service queues.
- **Response Publisher**: Publishes responses to the results queue with correlation IDs.

## Endpoints

### Gateway (Port 4000)
- `GET /` - Home endpoint returning status message.
- `GET /data` - Initiates scatter-gather pattern and returns aggregated responses.

### Service A (Port 4010)
- `GET /` - Home endpoint returning service status.
- `GET /data` - Returns dummy JSON data.

### Service B (Port 4020)
- `GET /` - Home endpoint returning service status.
- `GET /data` - Returns dummy JSON data.

## Message Flow

1. Client hits `/data` endpoint on gateway
2. Gateway generates unique correlation ID (UUID)
3. Gateway creates buffered channel and stores in sync.Map
4. Gateway publishes message to RabbitMQ exchange
5. Both services receive message from their respective queues
6. Services process message and publish responses to results queue
7. Gateway aggregator consumes responses and routes to appropriate channels
8. Gateway waits for responses (timeout: 10s, min responses: 2)
9. Gateway returns aggregated response to client

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
go mod tidy
go run main.go

# Terminal 2 - Start Service B
cd service_b
go mod tidy
go run main.go

# Terminal 3 - Start Gateway
cd gateway
go mod tidy
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

## RabbitMQ Management

Access RabbitMQ management interface at: http://localhost:15672
- Username: guest
- Password: guest

## Configuration

The system uses the following RabbitMQ configuration:
- **Exchange**: scatter-gather (fanout)
- **Service A Queue**: service_a_queue
- **Service B Queue**: service_b_queue
- **Results Queue**: results
- **Connection**: amqp://guest:guest@rabbitmq:5672/

## Stopping the System

```bash
# If using Docker Compose
docker-compose down

# If running manually, use Ctrl+C in each terminal
```

## Key Features

- **Correlation ID Tracking**: Each request gets a unique UUID for tracking responses.
- **Buffered Channels**: Uses Go's buffered channels for efficient response handling.
- **Timeout Handling**: Configurable timeout (10 seconds) for response collection.
- **Minimum Response Count**: Waits for at least 2 responses before returning.
- **Graceful Degradation**: Returns partial results if timeout occurs.
- **Concurrent Processing**: Services process messages independently and concurrently.
