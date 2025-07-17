# Database Migrations

## MySQL Setup

### 1. Install MySQL
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server

# macOS
brew install mysql
brew services start mysql
```

### 2. Create Database and User
```bash
# Connect to MySQL
mysql -u root -p

# Create database
CREATE DATABASE image_rag CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# Create user (optional)
CREATE USER 'imagerag'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON image_rag.* TO 'imagerag'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 3. Run Schema
```bash
mysql -u root -p image_rag < migrations/init.sql
```

## Milvus Setup

### 1. Install Milvus
```bash
# Using Docker
docker run -d --name milvus-standalone \
  -p 19530:19530 \
  -p 9091:9091 \
  -v $(pwd)/milvus/db:/var/lib/milvus/db \
  -v $(pwd)/milvus/conf:/var/lib/milvus/conf \
  -v $(pwd)/milvus/logs:/var/lib/milvus/logs \
  -v $(pwd)/milvus/wal:/var/lib/milvus/wal \
  milvusdb/milvus:v2.3.3

# Verify installation
docker logs milvus-standalone
```

### 2. Environment Variables
Create `.env` file in backend directory:
```
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=image_rag

# Milvus
MILVUS_HOST=localhost
MILVUS_PORT=19530
MILVUS_DATABASE=image_rag

# Doubao API
DOUBAO_API_KEY=your_api_key_here
DOUBAO_MODEL=doubao-embedding-vision-250615
DOUBAO_API_URL=https://ark.cn-beijing.volces.com/api/v3/embeddings

# Server
SERVER_PORT=8080
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE_MB=10
```

## Testing Connection

### 1. MySQL Test
```bash
mysql -u root -p -e "USE image_rag; SELECT * FROM records;"
```

### 2. Milvus Test
```bash
# Using Python (if installed)
python -c "from pymilvus import connections; connections.connect('default', host='localhost', port='19530')"
```

## Reset Database
```bash
# Reset MySQL
mysql -u root -p -e "DROP DATABASE image_rag; CREATE DATABASE image_rag CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p image_rag < migrations/init.sql

# Reset Milvus (if using Docker)
docker stop milvus-standalone
docker rm milvus-standalone
docker run -d --name milvus-standalone -p 19530:19530 -p 9091:9091 milvusdb/milvus:v2.3.3
```