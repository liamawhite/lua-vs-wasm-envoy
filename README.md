# Envoy Memory Usage Comparison: Lua vs WASM Filters

A performance comparison between Lua and Go WASM filters in Envoy, measuring memory usage under identical workloads.

## Results

**Memory Usage Comparison:**
- **Lua Filter**: ~17MB memory usage (lightweight)
- **WASM Plugin**: ~118MB memory usage (~7x overhead due to WASM runtime)

Monitor memory usage of the running services:

```bash
# Real-time monitoring
docker stats main-lua-envoy-1 main-wasm-v8-envoy-1 main-wasm-wamr-envoy-1

# Single snapshot
docker stats --no-stream main-lua-envoy-1 main-wasm-v8-envoy-1 main-wasm-wamr-envoy-1
```

```
# concurrency: 1
CONTAINER ID   NAME                     CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
180e29066d6d   main-lua-envoy-1         0.37%     19.79MiB / 15.6GiB   0.12%     2.59MB / 2.26MB   0B / 0B     10
42c641589a8b   main-wasm-v8-envoy-1     0.50%     117.7MiB / 15.6GiB   0.74%     2.61MB / 2.27MB   0B / 0B     25
2975d75ef974   main-wasm-wamr-envoy-1   0.42%     110.6MiB / 15.6GiB   0.69%     2.62MB / 2.27MB   0B / 0B     25

# concurrency: 2
CONTAINER ID   NAME                     CPU %     MEM USAGE / LIMIT    MEM %     NET I/O         BLOCK I/O   PIDS
090a1756129f   main-lua-envoy-1         0.51%     23.07MiB / 15.6GiB   0.14%     333kB / 290kB   0B / 0B     12
11009048820c   main-wasm-v8-envoy-1     0.42%     132.5MiB / 15.6GiB   0.83%     338kB / 292kB   0B / 0B     27
bc971b7cacba   main-wasm-wamr-envoy-1   0.38%     122.8MiB / 15.6GiB   0.77%     339kB / 293kB   0B / 0B     27

# concurrency: 4
CONTAINER ID   NAME                     CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
5c80df4594aa   main-lua-envoy-1         0.66%     22.63MiB / 15.6GiB   0.14%     78.5kB / 67.2kB   0B / 0B     16
feba0f0b2025   main-wasm-v8-envoy-1     0.54%     152MiB / 15.6GiB     0.95%     76.9kB / 65.5kB   0B / 0B     31
9976a50ec1bb   main-wasm-wamr-envoy-1   0.53%     151MiB / 15.6GiB     0.95%     79.4kB / 67.6kB   0B / 0B     31

# concurrency: 8
CONTAINER ID   NAME                     CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
4198d1847900   main-lua-envoy-1         0.84%     23.97MiB / 15.6GiB   0.15%     19.2kB / 15.4kB   0B / 0B     24
9e28ef777928   main-wasm-v8-envoy-1     0.75%     190.9MiB / 15.6GiB   1.20%     21.9kB / 17.6kB   0B / 0B     39
5b0a64004156   main-wasm-wamr-envoy-1   0.70%     189.2MiB / 15.6GiB   1.18%     22kB / 17.5kB     0B / 0B     39

# concurrency: 16
CONTAINER ID   NAME                     CPU %     MEM USAGE / LIMIT    MEM %     NET I/O           BLOCK I/O   PIDS
2706917ceea4   main-lua-envoy-1         0.82%     25.66MiB / 15.6GiB   0.16%     42.3kB / 35.7kB   0B / 0B     40
4c870b532cdc   main-wasm-v8-envoy-1     0.99%     291.8MiB / 15.6GiB   1.83%     40.8kB / 34kB     0B / 0B     55
3aac48975d8c   main-wasm-wamr-envoy-1   0.96%     277.9MiB / 15.6GiB   1.74%     43kB / 35.7kB     0B / 0B     55
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
