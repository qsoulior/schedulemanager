# Schedule Manager

Schedule Manager is app for 
1. downloading and parsing PDF schedules with specific layout,
2. providing these parsed schedules in JSON via API.

## Installation with Docker Compose

### For development
```
docker compose -f docker-compose.dev.yml --env-file <your_env_file> up -d --build
```

### For production
```
docker compose -f docker-compose.prod.yml --env-file <your_env_file> up -d --build
```

## Setting environment variables

### Required Moodle credentials
```
MOODLE_USERNAME=user
MOODLE_PASSWORD=password
MOODLE_ROOT_URL=https://example.com
MOODLE_COURSE_ID=12345
```

### MongoDB credentials
#### a. For development: Optional MongoDB user and password
> By default user is *user1*, password is *test1*
```
MONGO_USER=user1
MONGO_PASSWORD=test1
```

#### b. For production: Required MongoDB connection string
```
MONGODB_CONNSTRING=mongodb+srv://username:password@mongodb0.example.com:27017
```

### Optional web API access token

```
WEB_API_TOKEN=token
```