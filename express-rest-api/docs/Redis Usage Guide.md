# Redis Usage Guide for the Express REST API Project

This document provides a comprehensive guide to using Redis in the `express-rest-api` project, covering connection methods, commands, configuration, and troubleshooting.

## 1. Connecting to Redis

### Option 1: Using Globally Installed `redis-cli`

If `redis-cli` is installed globally on your system, you can connect directly:

```bash
redis-cli -h localhost -p 6379
```

### Option 2: Using Docker (No Global `redis-cli`)

If `redis-cli` is not installed globally, you can connect to Redis via Docker:

```bash
# Connect to a running Redis container
docker exec -it redis redis-cli

# Or, if the container has a different name
docker exec -it [redis_container_name] redis-cli
```

## 2. Redis Commands

### 2.1. Basic Commands

```bash
# Test connection
PING
# Response: PONG (if connection is successful)

# List all keys
KEYS *

# View server information
INFO

# View memory statistics
INFO memory

# Clear all data in the current database
FLUSHDB

# Clear all data in all databases
FLUSHALL

# Monitor Redis operations in real-time
MONITOR
# (Press Ctrl+C to exit MONITOR mode)
```

### 2.2. String Operations

```bash
# Set a value
SET key value

# Get a value
GET key

# Delete a key
DEL key

# Check if a key exists
EXISTS key

# Set expiration time (in seconds)
EXPIRE key seconds
```

### 2.3. Hash Map Operations (Used in the Project)

```bash
# Set a field in a hash
HSET hash field value

# Get the value of a field in a hash
HGET hash field

# Get all fields and values in a hash
HGETALL hash

# Delete a field in a hash
HDEL hash field

# Check if a field exists in a hash
HEXISTS hash field
```

In the `express-rest-api` project, Hash Maps are used with the `participants` key for caching:

```bash
# View all participants in the cache
HGETALL participants

# View details of a specific participant (replace [id] with the actual ID)
HGET participants [id]
```

### 2.4. Bull Queue Operations (Used in the Project)

```bash
# View all keys related to Bull queues
KEYS bull:*

# View pending jobs in the queue
LRANGE bull:participant-cache-refresh:wait 0 -1

# View completed jobs
LRANGE bull:participant-cache-refresh:completed 0 -1

# View failed jobs
LRANGE bull:participant-cache-refresh:failed 0 -1
```

## 3. Redis Configuration in the Project

The `express-rest-api` project uses Redis for two primary purposes:

1. **Data Caching**: Stores `participants` data to reduce database load.
2. **Job Queue**: Uses Bull to schedule periodic cache refreshes.

### Environment Variables (.env)

```env
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Docker Configuration (docker-compose.yml)

```yaml
redis:
  image: redis/redis-stack-server:latest
  container_name: redis
  ports:
    - "6379:6379"
```

## 4. Starting Redis with Docker

```bash
# Navigate to the project directory
cd /Volumes/SSD-English/software-engineer/learning-to-learn/back-end/2.\ general-course/express-rest-api/

# Start the Redis container
docker-compose up -d redis

# Verify the Redis container is running
docker ps | grep redis
```

## 5. GUI Tools for Redis (Alternative to Command Line)

If you prefer a graphical interface over the command line, consider these tools:

- **Redis Desktop Manager**: https://rdm.dev/
- **Redis Insight**: https://redis.com/redis-enterprise/redis-insight/
- **Medis**: https://getmedis.com/

**Connection Settings**:
- Host: `localhost`
- Port: `6379`
- Password: (none)

## 6. Troubleshooting

### Redis Connection Issues

1. **Check if the Redis container is running**:

   ```bash
   docker ps | grep redis
   ```

2. **Restart the Redis container if needed**:

   ```bash
   docker-compose restart redis
   ```

3. **Inspect Redis container logs**:

   ```bash
   docker logs redis
   ```

4. **Check if port 6379 is in use by another application**:

   ```bash
   lsof -i :6379
   ```

## Conclusion

This guide provides the necessary information to effectively use Redis in the `express-rest-api` project. For further assistance, refer to the official Redis documentation or the project’s configuration files.