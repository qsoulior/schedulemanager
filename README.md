# Schedule Manager

Schedule Manager is app for 
1. downloading and parsing PDF schedules with specific layout,
2. providing these parsed schedules in JSON via web API.

## Run with Docker Compose

### For development
There must be file `configs/docker.dev.json` before running.
> Default MongoDB connection string is `mongodb://user1:test1@host.docker.internal:27017`.
```
docker compose -f docker-compose.dev.yml up -d
```

### For production
There must be file `configs/docker.prod.json` before running.
```
docker compose -f docker-compose.prod.yml up -d
```

## Configuration
JSON configuration files are stored in folder `configs`.
Config path must be stored in `CONFIG_PATH` environment variable.

`configs/example.json`
```json
{
  "server": {
    "port": 3000,
    "allowed_origins": "*" 
  },
  "mongo": {
    "uri": "mongodb://user1:test1@example.host:27017/"
  },
  "moodle": {
    "host": "https://example.com",
    "username": "",
    "password": "",
    "course_id": 0
  }
}
```
