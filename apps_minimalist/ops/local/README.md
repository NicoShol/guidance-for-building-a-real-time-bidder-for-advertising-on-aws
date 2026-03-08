# Local Testing Scripts

Scripts for testing the minimalist bidder locally.

## Prerequisites

### Install vegeta

**macOS:**
```bash
brew install vegeta
```

**Linux/Other:**
```bash
go install github.com/tsenart/vegeta@latest
```

**Check installation:**
```bash
vegeta -version
```

## Scripts

### test-server
Builds and runs the bidder server locally.

```bash
./ops/local/test-server
```

Configuration is loaded from [test-env](test-env):
- Server: `localhost:8000`
- Endpoint: `/bidrequest`
- Timeout: `100ms`
- Log level: `DEBUG`

### test-fixtures
Sends a single bid request from a fixture file using curl.

```bash
./ops/local/test-fixtures banner-site
```

Available fixtures (in [../../fixtures/](../../fixtures/)):
- `banner-site` - Banner ad on website
- `banner-app` - Banner ad in mobile app
- `video-site` - Video ad on website
- `audio-site` - Audio ad on website

### run-load-test
Runs a basic load test using vegeta with static fixture data.

**Basic usage:**
```bash
./ops/local/run-load-test
```

**Custom configuration:**
```bash
# Set via environment variables
LOAD_TEST_RATE=200 \
LOAD_TEST_DURATION=60s \
LOAD_TEST_WORKERS=8 \
LOAD_TEST_FIXTURE=video-site \
./ops/local/run-load-test
```

**Configuration options:**
- `LOAD_TEST_RATE` - Requests per second (default: 100)
- `LOAD_TEST_DURATION` - Test duration (default: 30s)
- `LOAD_TEST_WORKERS` - Number of workers (default: 4)
- `LOAD_TEST_TIMEOUT` - Request timeout (default: 1s)
- `LOAD_TEST_FIXTURE` - Fixture to use (default: banner-site)

**Output:**
- Summary statistics (latency, success rate, throughput)
- Latency histogram
- Detailed metrics in JSON format
- Results automatically saved to `ops/local/testing_output/`:
  - `vegeta-results-TIMESTAMP.bin` - Binary results
  - `report-TIMESTAMP.json` - JSON metrics
  - `plot-TIMESTAMP.html` - Interactive HTML plot

**View HTML plot:**
```bash
# Plot is automatically generated
open ops/local/testing_output/plot-*.html
```

### run-load-test-dynamic
Runs a load test with dynamically generated request IDs (similar to the production load generator).

**Usage:**
```bash
./ops/local/run-load-test-dynamic
```

This script:
- Generates unique IDs for each request (`id`, `user.id`, `source.tid`, etc.)
- Simulates more realistic traffic patterns
- Uses the same configuration options as `run-load-test`
- Results saved with `-dynamic-` suffix in filenames

**Note:** Requires Go to be installed (to build the request generator).

## Example Workflow

```bash
# Terminal 1: Start the bidder
cd apps_minimalist
./ops/local/test-server

# Terminal 2: Run a quick test
./ops/local/test-fixtures banner-site

# Terminal 3: Run load test
./ops/local/run-load-test

# Or with custom settings
LOAD_TEST_RATE=500 LOAD_TEST_DURATION=2m ./ops/local/run-load-test

# View HTML plot (automatically generated)
open ops/local/testing_output/plot-*.html
```

## Interpreting Results

### Key Metrics

**Latency percentiles:**
- `p50` (median) - Half of requests complete faster than this
- `p95` - 95% of requests complete faster than this
- `p99` - 99% of requests complete faster than this
- `max` - Slowest request

**Success rate:**
- Should be close to 100% for healthy service
- Drops indicate errors or timeouts

**Throughput:**
- Actual requests/second achieved
- Compare to target rate to check if server is keeping up

### Example Output

```
Requests      [total, rate, throughput]         3000, 100.03, 99.98
Duration      [total, attack, wait]             30.007s, 29.990s, 17.234ms
Latencies     [min, mean, 50, 90, 95, 99, max]  5.123ms, 12.456ms, 11.234ms, 18.567ms, 22.345ms, 35.678ms, 78.901ms
Bytes In      [total, mean]                     450000, 150.00
Bytes Out     [total, mean]                     3600000, 1200.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:3000
```

**What to look for:**
- ✅ Success rate = 100%
- ✅ Throughput ≈ Rate (server keeping up)
- ✅ p95 < timeout threshold
- ⚠️ If throughput < rate: server is bottlenecked
- ⚠️ If p99 is high: some requests are slow (investigate)

## Troubleshooting

**"vegeta: command not found"**
- Install vegeta (see Prerequisites)

**"Server may not be running"**
- Start the server in another terminal: `./ops/local/test-server`

**"Fixture not found"**
- Check available fixtures: `ls ../../fixtures/`
- Use fixture name without `.json` extension

**High error rate**
- Check server logs for errors
- Increase timeout: `LOAD_TEST_TIMEOUT=5s`
- Reduce rate: `LOAD_TEST_RATE=50`

**Results show 0 req/s**
- Server might be rejecting all requests
- Check server logs for error messages
- Test with single request first: `./ops/local/test-fixtures banner-site`
