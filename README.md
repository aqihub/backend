# Go Tropic Thunder

A decentralized environmental data collection and incentivization platform built with Go, IPFS, Redis, and Ethereum blockchain technology.

## Overview

Go Tropic Thunder is a backend service that enables IoT devices to submit environmental data (temperature, humidity, air quality, GPS coordinates) which gets stored on IPFS (InterPlanetary File System) with metadata managed in Redis. The platform includes a token-based incentivization system that rewards data contributors with ERC20 tokens on the Avalanche network.

## Features

- **Environmental Data Collection**: Accept and store environmental sensor data including:
    - GPS coordinates (latitude/longitude)
    - Temperature (Celsius)
    - Humidity levels
    - TVOC (Total Volatile Organic Compounds) in PPB
    - eCO2 levels in PPM
    - Air Quality Index (AQI)
- **IPFS Storage**: Decentralized storage of environmental data documents
- **Redis Metadata Management**: Fast metadata indexing and retrieval
- **Blockchain Incentivization**: Automatic ERC20 token transfers to reward data contributors
- **RESTful API**: Clean HTTP endpoints for data submission and retrieval
- **CORS Support**: Cross-origin resource sharing for web applications
- **Dockerized Deployment**: Complete containerized setup with all dependencies

## Tech Stack

- **Backend**: Go 1.24
- **Storage**: IPFS (InterPlanetary File System)
- **Database**: Redis
- **Blockchain**: Ethereum/Avalanche (ERC20 tokens)
- **Containerization**: Docker
- **HTTP Router**: Gorilla Mux
- **IPFS Client**: go-ipfs-api
- **Ethereum Client**: go-ethereum

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for local development)
- Access to Ethereum/Avalanche RPC endpoint (Alchemy recommended)
- ERC20 token contract deployed

## Environment Variables

Create a `.env` file based on `.env.example`:

```bash
IPFS_DATABASE=your_database_name
REDIS_PORT=localhost:6379

# Blockchain Configuration
ALCHEMY_API_KEY=your_alchemy_api_key
ALCHEMY_ETH_NETWORK=https://api.avax-test.network/ext/bc/C/rpc
SIGNER_PRIVATE_KEY=your_private_key_without_0x_prefix
AQI_ERC20_ADDRESS=your_erc20_token_contract_address
```

## Installation & Setup

### Using Docker (Recommended)

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd go-tropic-thunder
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Build and run with Docker:**
   ```bash
   docker build -t go-tropic-thunder .
   docker run -p 3000:3000 -p 5001:5001 -p 4001:4001 -p 6379:6379 --env-file .env go-tropic-thunder
   ```

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start Redis and IPFS separately:**
   ```bash
   # Start Redis
   redis-server --port 6379

   # Start IPFS daemon
   ipfs daemon --enable-pubsub-experiment
   ```

3. **Run the application:**
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

### POST /insert
Submit environmental data from IoT devices.

**Request Body:**
```json
{
  "device_id": "sensor_001",
  "gps_lat": 40.7128,
  "gps_lng": -74.0060,
  "timestamp": 1640995200,
  "temp_cel": 22.5,
  "humidity": 65.0,
  "tvoc_ppb": 150.0,
  "eco2_ppm": 400.0,
  "aqi": 50,
  "is_public": true
}
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "cid": "QmXxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "txHash": "0xabcdef..."
  }
}
```

### GET /select
Retrieve stored environmental data.

**Query Parameters:**
- `collection_name`: Get all documents for a specific device/collection
- `document_id`: Get specific document by CID

**Examples:**
```bash
# Get all collections
curl http://localhost:3000/select

# Get documents for a specific device
curl http://localhost:3000/select?collection_name=sensor_001

# Get specific document by CID
curl http://localhost:3000/select?document_id=QmXxx...
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   IoT Device    │───▶│   Go Backend    │───▶│      IPFS       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐    ┌──────────────────┐
                       │      Redis      │    │   Blockchain     │
                       │   (Metadata)    │    │ (Incentivization)│
                       └─────────────────┘    └──────────────────┘
```

### Data Flow

1. **Data Submission**: IoT devices submit environmental data via HTTP POST
2. **IPFS Storage**: Raw data gets stored on IPFS, returning a Content ID (CID)
3. **Metadata Indexing**: Redis stores the mapping between device IDs and their CIDs
4. **Token Incentivization**: Successful submissions trigger automatic ERC20 token transfers
5. **Data Retrieval**: API endpoints allow querying stored data by device or CID

## Project Structure

```
go-tropic-thunder/
├── cmd/
│   └── main.go                 # Application entry point
├── pkg/
│   ├── api/
│   │   └── handler.go          # HTTP request handlers
│   ├── db/
│   │   └── metadata.go         # Redis metadata management
│   ├── incentivizer/
│   │   ├── constants.go        # Blockchain constants
│   │   └── incentivizer.go     # Token transfer logic
│   ├── models/
│   │   └── models.go           # Data structures
│   ├── routes/
│   │   └── routes.go           # HTTP routing and middleware
│   └── storage/
│       └── ipfs.go             # IPFS client wrapper
├── Dockerfile                  # Docker configuration
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── .env.example               # Environment variables template
└── .gitignore                 # Git ignore rules
```

## Docker Services

The Docker container runs multiple services:
- **Redis**: Port 6379 - Metadata storage
- **IPFS**: Ports 4001 (swarm), 5001 (API), 8080 (gateway)
- **Go Application**: Port 3000 - Main API server

## Health Check

The application includes a health check endpoint accessible at:
```
GET http://localhost:3000/health
```

## Token Incentivization

The platform automatically rewards data contributors with ERC20 tokens:
- **Default Amount**: 420 tokens per successful submission
- **Target Address**: Configured via `TransferAccountAddress` constant
- **Network**: Avalanche Testnet (configurable via environment variables)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## Support
Contributed by [MadhuS](https://github.com/MadhuS-1605) with ❤️
For support or questions, please open an issue on the GitHub repository.