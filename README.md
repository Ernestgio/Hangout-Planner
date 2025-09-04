# Hangout-Planner

Hangout Planner

## Project Tech Stack

- Go
- Echo
- GORM
- MySQL
- Docker

## Run Local Development

### Hangout Service

Execute

```
docker build -t hangout . && docker run --rm --env-file .env -p 9000:9000 hangout && docker rmi hangout
```
