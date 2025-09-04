# Hangout-Planner

Hangout Planner

## Project Tech Stack

- Go
- Echo
- GORM
- MySQL
- Docker

## Run Local Development

Execute

```
cd cmd/Hangout && docker build -t hangout . && docker run --rm --env-file .env -p 9000:9000 hangout
```
