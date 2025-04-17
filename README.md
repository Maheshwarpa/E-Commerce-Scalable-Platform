# SKHT E-Commerce Site - README

## Project Overview

SKHT is an e-commerce platform that allows users to log in, browse products, and create accounts. The backend architecture involves several databases (PostgreSQL and Cassandra) for handling various types of data such as user information, product listings, and orders. The system also utilizes a REST API for user authentication and product management, with a GRPC server handling communication between services.

qvmh uaqi ehqs vjnj

### Technology Stack

- Backend: Go (Golang)
- Databases:
- PostgreSQL (for user and product data)
- Cassandra (for order-related data)
- Messaging System: Kafka
- Search Engine: Elasticsearch
- Logging: Custom logger with structured logging
- Containerization: Docker and Docker Compose
- APIs: REST and GRPC
- Service Architecture: Microservices

### Project Structure

````
.
├── docker-compose.yml       # Docker Compose file to bring up all services
├── main.go                  # Main entry point for the application
├── go.mod                   # Go modules dependencies
├── go.sum                   # Go modules checksum
├── module/
│   ├── API/                 # REST API package for user interaction
│   ├── Cassandra/           # Package for interacting with Cassandra DB
│   ├── Database/            # Package for interacting with PostgreSQL DB
│   ├── ProductService/      # Service to manage products
│   ├── UserService/         # Service to handle user management
│   ├── logger/              # Custom logger package
└── logs/                    # Log directory for storing application logs

````

### Setup and Installation

1. Clone the Repository

First, clone this repository to your local machine
````
git clone https://github.com/your-repo/skht-ecommerce.git
cd skht-ecommerce
````

2. Docker Setup

Ensure that Docker is installed on your machine and Docker Compose is available. The docker-compose.yml file contains configuration for running the following services:

- PostgreSQL (for storing user and product information)
- Cassandra (for storing order-related information)
- Kafka (message queue for handling data streams)
- Zookeeper (for Kafka coordination)
- Redis (caching layer)
- Elasticsearch (for product search)
- Logstash (for log aggregation)
- Kibana (for visualizing logs)

Start the services by running:
```
docker-compose up --build

```

### Application Setup

1. Initialize Databases
When the application starts, the following actions will be performed automatically:

#### PostgreSQL Database:

- Connect to the database.
- Create the following tables: Product, User, and Login.
- Load initial product data into the database.

#### Cassandra Database:

- Connect to Cassandra.
- Create the Order and FinalOrder tables for managing order-related data.
2. Running the Application
The application’s entry point is the main.go file. The program offers a simple console-based menu for users:

- Login (Existing users)
- New User (Create a new user)

To run the application:
  ```
  go run main.go
  ```

### Services and Functionality

1. PostgreSQL Database (Database Package)
- Connection: Uses pgx to connect to PostgreSQL.
- Tables: Creates tables for products, users, and login details.
- Data Loading: Loads product data into the Product table.
2. Cassandra Database (Cassandra Package)
- Connection: Connects to the Cassandra database for storing orders.
- Tables: Creates Order and FinalOrder tables to manage order-related data.
- CQL Queries: Initializes and populates Cassandra with the necessary data.
3. User Management (UserService Package)
- User Creation: Allows new users to be created with basic details (like name and email).
- Authentication: Handles user login functionality through the REST API.
4. Product Management (ProductService Package)
- Product Data: Loads a list of products into the database and makes it accessible via API calls.
5. API Server (API Package)
- Login Endpoint: Allows users to log in via a REST API.
- Authentication: Uses JWT or similar methods for user authentication.
6. Logging (logger Package)
- Logger Initialization: Configures the application’s logger to capture important events and errors.
- Log Output: Logs are stored in the logs/ directory.
7. GRPC Server
- Starts a GRPC server for internal communication between services.

### Configuration and Environment Variables

The following environment variables are used to configure the services:

#### PostgreSQL Configuration:
- POSTGRES_USER: Username for PostgreSQL.
- POSTGRES_PASSWORD: Password for PostgreSQL.
- POSTGRES_DB: The database name to connect to.
#### Cassandra Configuration:
- CASSANDRA_CLUSTER_NAME: The name of the Cassandra cluster.
- CASSANDRA_NUM_TOKENS: The number of tokens for Cassandra.
#### Kafka and Zookeeper Configuration:
- KAFKA_CFG_ZOOKEEPER_CONNECT: The Zookeeper address.
- KAFKA_CFG_LISTENERS: The Kafka listener configurations.
- KAFKA_CFG_ADVERTISED_LISTENERS: The advertised Kafka listeners.
  
### Testing

#### Unit Tests
Unit tests are available for each package. To run the tests:
```
go test ./...

```
### Logs and Monitoring

- Logs are written to the logs/ directory.
- Elasticsearch and Kibana are used for visualizing and searching through the logs.
- Logstash is used to aggregate and ship logs into Elasticsearch.

### Useful Docker Commands

- To stop the services:
```
docker-compose down

```
- To view logs of a specific service:
```
docker-compose logs <service-name>

```
