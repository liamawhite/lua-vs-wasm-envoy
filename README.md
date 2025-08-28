# Envoy Memory Usage Comparison: Lua vs WASM Filters

A performance comparison between Lua and Go WASM filters in Envoy, measuring memory usage under identical workloads.

## Results

**Memory Usage Comparison:**
- **Lua Filter**: ~17MB memory usage (lightweight)
- **WASM Plugin**: ~118MB memory usage (~7x overhead due to WASM runtime)

Monitor memory usage of the running services:

```bash
# Real-time monitoring
docker stats wasm-lua-test-lua-envoy-1 wasm-lua-test-wasm-envoy-1

# Single snapshot
docker stats --no-stream wasm-lua-test-lua-envoy-1 wasm-lua-test-wasm-envoy-1
```

```
# concurrency: 1
CONTAINER ID   NAME                CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
0977bbcaf325   main-lua-envoy-1    0.34%     19.48MiB / 15.6GiB   0.12%     31.1kB / 25.9kB   0B / 0B     10
59bdea6cc12b   main-wasm-envoy-1   0.38%     128.7MiB / 15.6GiB   0.81%     29.2kB / 24.4kB   0B / 0B     25

# concurrency: 2
CONTAINER ID   NAME                CPU %     MEM USAGE / LIMIT    MEM %     NET I/O         BLOCK I/O   PIDS
10559651531b   main-lua-envoy-1    0.41%     22.83MiB / 15.6GiB   0.14%     337kB / 293kB   0B / 0B     12
26c535bbd7d6   main-wasm-envoy-1   0.35%     132.5MiB / 15.6GiB   0.83%     335kB / 291kB   0B / 0B     27

# concurrency: 4
CONTAINER ID   NAME                CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
1a6dfea5b398   main-lua-envoy-1    0.49%     21.37MiB / 15.6GiB   0.13%     49kB / 41.6kB     0B / 0B     16
960467ffbc23   main-wasm-envoy-1   0.45%     157.8MiB / 15.6GiB   0.99%     46.6kB / 39.5kB   0B / 0B     31

# concurrency: 8
CONTAINER ID   NAME                CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
0cb7fecc051f   main-lua-envoy-1    0.63%     23.88MiB / 15.6GiB   0.15%     33kB / 27.7kB     0B / 0B     24
a03c26738b70   main-wasm-envoy-1   0.62%     198.9MiB / 15.6GiB   1.25%     30.6kB / 25.6kB   0B / 0B     39

# concurrency: 16
CONTAINER ID   NAME                CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
1df08c7c1a72   main-lua-envoy-1    0.89%     27.49MiB / 15.6GiB   0.17%     32.5kB / 27.2kB   0B / 0B     40
c42ec15f8ea6   main-wasm-envoy-1   0.80%     284.5MiB / 15.6GiB   1.78%     30.2kB / 25.1kB   0B / 0B     55
```

Additional memory monitoring options:
```bash
# Formatted output
docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}\t{{.MemPerc}}"

# Envoy-specific memory stats
curl -s localhost:9901/stats | grep memory  # Lua service
curl -s localhost:9902/stats | grep memory  # WASM service
```

## Test Architecture

- **Lua Service** (port 10000): Envoy with inline Lua filter
- **WASM Service** (port 10001): Envoy with Go-compiled WASM plugin
- **Load Generators**: curl-based containers generating 5 RPS to each service
- **Identical Logic**: Both filters implement the same coin flip algorithm for fair comparison

## Usage

Start all services and load generators:
```bash
docker-compose up --build
```

Test the services manually:
```bash
# Lua filter (port 10000) - coin flip returns 200 or 500
curl localhost:10000
# Expected: "Lua filter: heads = 200" or "Lua filter: tails = 500"

# WASM plugin (port 10001) - coin flip returns 200 or 500
curl localhost:10001
# Expected: "WASM filter: heads = 200" or "WASM filter: tails = 500"
```

## Monitoring

Admin interfaces:
- Lua service: http://localhost:9901
- WASM service: http://localhost:9902

Load generators run continuously with 5 RPS to each service. Check logs:
```bash
docker-compose logs lua-load-gen
docker-compose logs wasm-load-gen
```

## Technical Details

- **Go Version**: 1.24.4
- **Envoy Version**: 1.33
- **WASM SDK**: github.com/proxy-wasm/proxy-wasm-go-sdk v0.0.0-20250212164326-ab4161dcf924
- **Compilation**: Go 1.24.4 with WASI target and c-shared buildmode
- **Coin Flip Logic**: Deterministic hash of all request headers (50/50 chance of 200/500)

## Files
- `lua_simple_test.yaml` - Lua filter configuration
- `wasm_simple_test.yaml` - WASM plugin configuration  
- `simple_filter.go` - Go WASM plugin source code
- `go.mod` / `go.sum` - Go module dependencies
- `Dockerfile-wasm` - Multi-stage build for WASM plugin
- `docker-compose.yml` - Services and load generators
