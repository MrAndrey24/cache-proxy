
# Caching Proxy Server

command-line HTTP caching proxy server built in Go. This tool forwards requests to an origin server, caches the successful responses locally, and serves subsequent identical requests from the cache to improve performance and reduce origin load.

This project is based on the [roadmap.sh Caching Proxy](https://roadmap.sh/projects/caching-server) specification and implemented using Clean Architecture principles.


## Features

* **HTTP Request Forwarding:** Transparently proxies incoming HTTP requests to a designated origin server.
* **Local Disk Caching:** Stores successful `GET` responses (headers, status codes, and body) in a local `cache_store.json` file.
* **Cache Status Headers:** Injects an `X-Cache` header (`HIT` or `MISS`) into the HTTP response.
* **Cache Management:** Includes a CLI flag to wipe the existing cache.
* **Clean Architecture:** Organized into Domain, Data, Service, and Handler layers for maintainability.


## Installation

Ensure you have [Go](https://go.dev/) installed. Clone the repository and build the binary:

```bash
# Navigate to the project root
cd caching-proxy

# Build the executable
go build -o caching-proxy main.go

```
## Usage/Examples

```bash
./caching-proxy --port 3000 --origin http://dummyjson.com
```

If you prefer to run it without building, use: 
```bash
go run main.go --port 3000 --origin http://dummyjson.com
```

## Clear the Cache

To delete the saved cache file and start fresh:

```bash
./caching-proxy --clear-cache
```
Or 

```bash
go run . --clear-cache
```


## Testing the Proxy

Once the server is running, you can test it using curl from a new terminal window.

- First Request (Cache MISS)

```bash
curl -i http://localhost:3000/products
```

Notice the header X-Cache: MISS in the response.


- Second Request (Cache HIT)

Run the exact same command again. The proxy will intercept the request and serve the data instantly from the local cache file.

```bash
curl -i http://localhost:3000/products
```

Notice the header X-Cache: HIT in the response.