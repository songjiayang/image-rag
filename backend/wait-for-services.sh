#!/bin/sh

# Debug environment
echo "Environment variables:"
echo "MYSQL_HOST: ${MYSQL_HOST:-mysql}"
echo "MYSQL_PORT: ${MYSQL_PORT:-3306}"
echo "MILVUS_HOST: ${MILVUS_HOST:-milvus}"
echo "MILVUS_PORT: ${MILVUS_PORT:-19530}"

# Function to check if a host:port is reachable
check_host() {
  local host=$1
  local port=$2
  local service=$3
  
  echo "Checking ${service} at ${host}:${port}..."
  
  # Try to resolve hostname
  if ! nslookup ${host} > /dev/null 2>&1; then
    echo "Warning: Cannot resolve hostname ${host}"
    return 1
  fi
  
  # Try to connect
  timeout 5 nc -z ${host} ${port}
  return $?
}

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
for i in $(seq 1 30); do
  if check_host ${MYSQL_HOST:-mysql} ${MYSQL_PORT:-3306} "MySQL"; then
    echo "MySQL is ready!"
    break
  fi
  echo "MySQL not ready, retrying in 2 seconds... (attempt ${i}/30)"
  sleep 2
done

# Wait for Milvus to be ready
echo "Waiting for Milvus to be ready..."
for i in $(seq 1 30); do
  if check_host ${MILVUS_HOST:-milvus} ${MILVUS_PORT:-19530} "Milvus"; then
    echo "Milvus is ready!"
    break
  fi
  echo "Milvus not ready, retrying in 2 seconds... (attempt ${i}/30)"
  sleep 2
done

# Additional wait for Milvus to fully initialize
echo "Waiting for Milvus to fully initialize..."
sleep 5

echo "All services are ready! Starting backend..."
exec "$@"