# go rate limiter

## Objective  
Develop a **rate limiter** in Go that can be configured to limit the maximum number of requests per second based on a specific **IP address** or an **access token**.  

## Description  
The goal of this challenge is to create a rate limiter in Go that can be used to control the request traffic to a web service. The rate limiter should be able to limit the number of requests based on two criteria:  

- **IP Address:** The rate limiter should restrict the number of requests received from a single IP address within a defined time interval.  
- **Access Token:** The rate limiter should also be able to limit requests based on a unique access token, allowing different expiration time limits for different tokens.  
  - The token must be provided in the request header in the following format:  
    ```text
    API_KEY: <TOKEN>
    ```
  - The access token limit settings should **override** the IP-based limits.  
    - Example: If the IP limit is **10 req/s** and a specific token limit is **100 req/s**, the rate limiter should use the **token-based** limit.  

## Requirements  

- The rate limiter must function as **middleware** that can be injected into the web server.  
- The rate limiter should allow configuring the **maximum number of requests** permitted per second.  
- The rate limiter should have an option to define a **block time** for an IP or token when the request limit is exceeded.  
- Configuration settings should be defined via **environment variables** or a `.env` file in the root directory.  
- The rate limiter should support **both IP-based and token-based** request limiting.  
- The system must respond appropriately when the limit is exceeded:  
  - **HTTP Status Code:** `429 Too Many Requests`  
  - **Message:** `"You have reached the maximum number of requests or actions allowed within a certain time frame."`  
- The rate-limiting logic should be **separated** from the middleware.  

## Examples  
### **IP-Based Limiting**  
If the rate limiter is configured to allow a maximum of **5 requests per second per IP**, and the IP `192.168.1.1` sends **6 requests in one second**, the **6th request should be blocked**.  

### **Token-Based Limiting**  
If a token `abc123` has a limit of **10 requests per second** and sends **11 requests within that timeframe**, the **11th request should be blocked**.  

### **Blocking Period**  
In both cases above, additional requests can only be made **after the expiration time has elapsed**.  
- Example: If the expiration time is **5 minutes**, the blocked IP or token will be able to send requests **only after 5 minutes**.  

---

## Local environment setup (Docker)
To run the project locally, you need to have Docker installed on your machine.

Create a `.env` file from the `.env.example` file in the root directory of the project.
```bash
cp .env.example .env
```

Run the following command to start the project:
```bash
make run
```

To remove the containers, run:
```bash
make down
```

## Rate limiter feature
The rate limiter is implemented as a middleware that can be injected into the web server. The rate limiter itself needs to be configured with the following environment variables:
- **RATE_LIMIT_IP_MAX_REQUEST**: Maximum number of requests per second per IP address. Default: 5
- **RATE_LIMIT_IP_BLOCK_TIME**: Block time for IP address in seconds. Default: 300 seconds (5 minutes)
- **RATE_LIMIT_TOKEN_MAX_REQUEST**: Maximum number of requests per second per token. Default: 10
- **RATE_LIMIT_TOKEN_BLOCK_TIME**: Block time for token in seconds. Default: 300 seconds (5 minutes)

Custom tokens can be set in the `customTokenBlockDuration` function with a desired expiration time.
 - the token `a6b3fdef-c107-4970-8ecc-94817ed5968c` has a block time of 30 seconds.


The rate limiter is designed based on the **Ports and Adapters** architecture. In this implementation, **Redis** serves as the adapter. The adapter can be replaced with any other database or in-memory storage that implements the port interface.

## Sending requests
After start the project, you can send requests using IP or token in the header.

Curl using token:
```bash
curl -X GET http://localhost:8080 \
		-H "API_KEY: 2c02b5ce-04d0-4c75-9810-c3e75c397956"
```

Curl using local IP:
```bash
curl -X GET http://localhost:8080
```

Curl using custom IP:
```bash
curl -X GET http://localhost:8080 \
		-H "X-Forwarded-For: 192.168.0.1"
```

You can use the predefined requests in the `Makefile` and `api.http` files to test the rate limiter.
- **Makefile**: `make request-token-2-min`
- **api.http**: you will need the `REST Client` extension.

## Integration test
To run the integration tests, ensure that Docker and Go are installed on your machine. The tests use the `testcontainers` library to create a Redis container and execute the tests through HTTP requests.

To run the integration test, run the following command in the root directory of the project:
```bash
make integration-test
```

> ⚠️ The integration tests include scenarios where IPs and tokens are blocked and then unblocked. As a result, the integration test may take some time to complete.

## Load testing
To run the load tests, ensure that Docker is installed on your machine. We will use the `hey` tool to perform the load testing. The tests will be executed against the rate limiter to check its performance under high load.

Start the project using Docker:
```bash
make run
```

In another terminal, run any command with the prefix `load-test-*` in the `Makefile` to perform the load test. For example:
```bash
make load-test-ip
```

**Hey tool**: the `-n` flag specifies the total number of requests to be sent and the `-c` flag specifies the number of concurrent requests.
