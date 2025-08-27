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

Example output:
```
CONTAINER ID   NAME                         CPU %     MEM USAGE / LIMIT     MEM %     NET I/O         BLOCK I/O    PIDS
5088bd126756   wasm-lua-test-lua-envoy-1    1.09%     17.14MiB / 7.656GiB   0.22%     528kB / 460kB   4.1kB / 0B   24
0dbafc7343b2   wasm-lua-test-wasm-envoy-1   1.29%     118.5MiB / 7.656GiB   1.51%     291kB / 254kB   0B / 0B      31
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