# Connection Pool Health Check and Reconnection

This document describes the connection pool health check and automatic reconnection features implemented in the database package.

## Features

### 1. Automatic Health Checks

The database connection pool automatically performs health checks every 30 seconds to ensure the database connection is alive and responsive.

**Key Features:**
- Periodic health checks using `PingContext` with 5-second timeout
- Automatic reconnection on health check failure
- Detailed statistics tracking

### 2. Automatic Reconnection

When a health check fails, the system automatically attempts to reconnect to the database:

**Reconnection Process:**
1. Close the old database connection
2. Create a new connection with the same configuration
3. Reconfigure the connection pool settings
4. Verify the new connection with a ping
5. Update health status and statistics

### 3. Connection Pool Statistics

The system tracks comprehensive statistics about the connection pool:

**Basic Statistics:**
- `MaxOpenConnections`: Maximum number of open connections
- `OpenConnections`: Current number of open connections
- `InUse`: Number of connections currently in use
- `Idle`: Number of idle connections
- `WaitCount`: Total number of times waited for a connection
- `WaitDuration`: Total time spent waiting for connections
- `MaxIdleClosed`: Connections closed due to idle timeout
- `MaxLifetimeClosed`: Connections closed due to max lifetime

**Health Check Statistics:**
- `LastHealthCheck`: Timestamp of last health check
- `HealthCheckCount`: Total number of health checks performed
- `HealthCheckFailed`: Number of failed health checks
- `LastHealthStatus`: Result of last health check (true/false)

**Reconnection Statistics:**
- `ReconnectCount`: Number of reconnection attempts
- `LastReconnectTime`: Timestamp of last reconnection
- `LastReconnectError`: Error message from last reconnection (if any)

## Usage

### Starting Health Checks

Health checks are automatically started when you initialize the database:

```go
import "github.com/Kevin-Jii/tower-go/utils/database"

// Initialize database (health checks start automatically)
err := database.InitDB(dbConfig)
if err != nil {
    log.Fatal(err)
}
```

### Manually Starting/Stopping Health Checks

You can manually control health checks if needed:

```go
// Start health checks
database.StartHealthCheck()

// Stop health checks
database.StopHealthCheck()
```

### Checking Connection Health

```go
// Check if the database connection is healthy
if database.IsHealthy() {
    fmt.Println("Database connection is healthy")
} else {
    fmt.Println("Database connection is unhealthy")
}
```

### Getting Connection Pool Statistics

```go
// Get detailed connection pool statistics
stats, err := database.GetConnectionPoolStats()
if err != nil {
    log.Printf("Error getting stats: %v", err)
    return
}

fmt.Printf("Open Connections: %d\n", stats.OpenConnections)
fmt.Printf("In Use: %d\n", stats.InUse)
fmt.Printf("Idle: %d\n", stats.Idle)
fmt.Printf("Health Check Count: %d\n", stats.HealthCheckCount)
fmt.Printf("Health Check Failed: %d\n", stats.HealthCheckFailed)
fmt.Printf("Reconnect Count: %d\n", stats.ReconnectCount)
fmt.Printf("Last Health Status: %v\n", stats.LastHealthStatus)
```

### Getting Raw SQL Statistics

```go
// Get raw sql.DBStats
rawStats, err := database.GetRawConnectionPoolStats()
if err != nil {
    log.Printf("Error getting raw stats: %v", err)
    return
}

fmt.Printf("Max Open Connections: %d\n", rawStats.MaxOpenConnections)
fmt.Printf("Wait Count: %d\n", rawStats.WaitCount)
fmt.Printf("Wait Duration: %v\n", rawStats.WaitDuration)
```

## Configuration

Health check behavior can be configured through the performance configuration:

```go
// In config/performance.go
type DatabasePerformanceConfig struct {
    MaxOpenConns    int           // Maximum open connections
    MaxIdleConns    int           // Maximum idle connections
    ConnMaxLifetime time.Duration // Connection max lifetime
    ConnMaxIdleTime time.Duration // Connection max idle time
    // ... other fields
}
```

Environment variables:
- `PERF_DB_MAX_OPEN_CONNS`: Maximum open connections (default: CPU count * 2)
- `PERF_DB_MAX_IDLE_CONNS`: Maximum idle connections (default: CPU count)
- `PERF_DB_CONN_MAX_LIFETIME_MINUTES`: Connection max lifetime in minutes (default: 60)
- `PERF_DB_CONN_MAX_IDLE_TIME_MINUTES`: Connection max idle time in minutes (default: 10)

## Implementation Details

### Health Check Interval

Health checks run every 30 seconds. This interval is hardcoded but can be made configurable if needed.

### Health Check Timeout

Each health check has a 5-second timeout. If the database doesn't respond within this time, the check is considered failed.

### Reconnection Strategy

The system uses a simple reconnection strategy:
1. On health check failure, immediately attempt reconnection
2. If reconnection fails, log the error and wait for the next health check
3. No exponential backoff or retry limits (relies on periodic health checks)

### Thread Safety

All statistics and connection management operations are protected by mutexes:
- `dbMutex`: Protects database instance during reconnection
- `statsLock`: Protects statistics updates

## Monitoring

The health check system logs important events:

**Info Level:**
- Connection pool health check started
- Connection pool health check stopped
- Database reconnected successfully

**Debug Level:**
- Connection pool statistics (on successful health check)

**Error Level:**
- Database health check failed
- Database reconnection failed

**Warning Level:**
- Attempting to reconnect to database

## Best Practices

1. **Monitor Statistics**: Regularly check connection pool statistics to identify issues
2. **Alert on Failures**: Set up alerts for repeated health check failures
3. **Tune Connection Pool**: Adjust pool size based on your workload
4. **Review Logs**: Check logs for reconnection events and errors

## Requirements Validation

This implementation validates the following requirements:

- **Requirement 5.3**: Connection pool queue - Requests wait when pool is exhausted
- **Requirement 5.4**: Unhealthy connection handling - Automatic reconnection on failure

## Testing

Run the tests with:

```bash
go test -v ./utils/database -run "TestConnectionPoolStats|TestIsHealthy|TestStartStopHealthCheck|TestPingWithContext|TestGetRawConnectionPoolStats"
```

All tests should pass, validating:
- Statistics tracking
- Health status checking
- Health check start/stop
- Context-aware ping
- Raw statistics retrieval
