# NGINX API Gateway

This component serves as the centralized entry point for the Hangout Planner microservices architecture, providing reverse proxying, SSL/TLS termination, and security hardening.

## Technical Overview

The infrastructure is built on NGINX to handle high-concurrency traffic with minimal resource overhead. It implements industry-standard best practices for security and performance.

### Key Features

- **Reverse Proxy & Path-based Routing**: Centralized routing logic using URL rewriting to decouple external API paths from internal service structures.
- **SSL/TLS Termination**: Secure HTTPS communication using modern TLS 1.2/1.3 protocols.
- **HTTP/2 Support**: Optimized binary protocol for reduced latency and improved multiplexing.
- **Security Hardening**: Implementation of HSTS, X-Frame-Options, X-Content-Type-Options, and Referrer-Policy headers.
- **Edge Compression**: Global Gzip compression to reduce bandwidth consumption and improve time-to-first-byte (TTFB).
- **Infrastructure as Code**: Fully containerized setup integrated with Docker Compose for consistent deployment across environments.

## Architecture

External requests are received on port 443 (HTTPS) and routed internally to the appropriate service over the Docker bridge network.

```
Client (HTTPS) -> NGINX (SSL Termination) -> Internal Service (HTTP)
```

### Routing Table

| External Path                                | Internal Upstream      | Description         |
| -------------------------------------------- | ---------------------- | ------------------- |
| `https://localhost/`                         | Static JSON            | Gateway Health/Info |
| `https://localhost/healthz`                  | `hangout:9000/healthz` | Global Health Check |
| `https://localhost/rp-api/hangout-service/*` | `hangout:9000/*`       | Hangout Service API |

## Configuration Details

### Security Headers

The gateway injects the following security headers into all responses:

- `Strict-Transport-Security`: Enforces HTTPS for 1 year.
- `X-Frame-Options`: Prevents clickjacking via `SAMEORIGIN` policy.
- `X-Content-Type-Options`: Disables MIME-type sniffing.
- `Referrer-Policy`: Implements `strict-origin-when-cross-origin`.

### Performance Tuning

- **Gzip**: Enabled for all text-based MIME types with a minimum length threshold of 512 bytes.
- **Keepalive**: Optimized connection pooling for upstream services.
- **Buffering**: Tuned proxy buffers (8k/64k) to handle typical JSON API payloads efficiently.

## Local Development Setup

For local development, the gateway uses `mkcert` to provide locally-trusted SSL certificates, ensuring the development environment mirrors production security.

### Certificate Generation

If certificates need to be regenerated:

```powershell
cd components/nginx/ssl
mkcert localhost 127.0.0.1 ::1
# mkcert generates localhost+2.pem and localhost+2-key.pem
# Rename them to match nginx configuration:
mv localhost+2.pem localhost.pem
mv localhost+2-key.pem localhost-key.pem
```

## Operational Commands

Commands are executed from the `services/hangout/` directory via the provided Makefile.

### Service Management

- `make up`: Initialize all services including the API Gateway.
- `make down`: Teardown the entire infrastructure.
- `make nginx-reload`: Hot-reload NGINX configuration without downtime.
- `make nginx-test`: Validate NGINX configuration syntax.

### Monitoring

- `make nginx-logs`: Stream real-time access and error logs.

## Directory Structure

```
nginx/
├── conf/
│   ├── nginx.conf          # Global NGINX configuration
│   └── conf.d/
│       └── hangout.conf    # Service-specific routing rules
├── ssl/                    # SSL/TLS Certificates (mkcert)
└── logs/                   # Access and Error logs
```
