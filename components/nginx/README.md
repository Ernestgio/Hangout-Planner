# NGINX API Gateway Setup - Phase 1: HTTP Only

This directory contains the NGINX reverse proxy configuration for the Hangout Planner API Gateway.

## 📁 Directory Structure

```
nginx/
├── conf/
│   ├── nginx.conf              # Main NGINX configuration
│   └── conf.d/
│       └── hangout.conf        # Hangout service reverse proxy config
├── logs/                       # NGINX access and error logs
├── ssl/                        # SSL certificates (for Phase 2)
└── README.md                   # This file
```

## 🚀 Phase 1: HTTP-Only Setup (Current)

### Architecture

```
http://localhost/rp-api/hangout-service/auth/signin
    ↓
NGINX (Port 80) - URL Rewrite
    ↓
http://hangout:9000/auth/signin (internal Docker network)
```

### Quick Start

From project root (`D:\ERNEST\Hangout-Planner`):

```powershell
# Start all services (MySQL + Hangout + NGINX)
docker-compose up -d

# Check service health
docker ps

# View NGINX logs
docker logs hangout-nginx -f
```

Or from `services/hangout/` directory:

```powershell
# Start all services with NGINX
make up

# View NGINX logs
make nginx-logs

# Test NGINX config
make nginx-test

# Reload NGINX after config changes
make nginx-reload

# Stop all services
make down-nginx
```

### Test Endpoints

```powershell
# Gateway info
curl http://localhost/

# Health check
curl http://localhost/healthz

# API endpoints (through reverse proxy)
curl http://localhost/rp-api/hangout-service/healthz
curl http://localhost/rp-api/hangout-service/swagger/index.html

# Sign up
curl -X POST http://localhost/rp-api/hangout-service/auth/signup `
  -H "Content-Type: application/json" `
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

### URL Routing

| External URL                                | Internal URL                  | Description    |
| ------------------------------------------- | ----------------------------- | -------------- |
| `http://localhost/`                         | -                             | Gateway info   |
| `http://localhost/healthz`                  | `http://hangout:9000/healthz` | Health check   |
| `http://localhost/rp-api/hangout-service/*` | `http://hangout:9000/*`       | All API routes |

**Example:**

- External: `http://localhost/rp-api/hangout-service/auth/signin`
- NGINX strips: `/rp-api/hangout-service`
- Backend receives: `/auth/signin`

### Configuration Details

**Key Features:**

- ✅ Path-based routing with URL rewriting
- ✅ Gzip compression at edge (removed from app)
- ✅ Proxy header forwarding (`X-Real-IP`, `X-Forwarded-For`, etc.)
- ✅ Security headers (`X-Frame-Options`, `X-Content-Type-Options`, etc.)
- ✅ Connection keepalive and buffering optimization
- ✅ Health check monitoring

**NGINX performs:**

1. Receives request on port 80
2. Strips `/rp-api/hangout-service` prefix
3. Forwards to `hangout:9000` with clean path
4. Adds proxy headers for client information
5. Compresses response (gzip)

### Makefile Commands

From `services/hangout/` directory:

```powershell
# Start/Stop
make up-nginx       # Start all services with NGINX
make down-nginx     # Stop all services

# NGINX Management
make nginx-reload   # Reload NGINX config (no downtime)
make nginx-test     # Test NGINX configuration validity
make nginx-logs     # Follow NGINX logs in real-time

# Development
make air            # Run with live reload (without NGINX)
make test           # Run tests
make swag           # Regenerate Swagger docs
```

### Troubleshooting

**NGINX won't start:**

```powershell
# Test configuration
make nginx-test

# Check logs
make nginx-logs

# Verify port 80 is available
netstat -ano | findstr :80
```

**502 Bad Gateway:**

- Hangout service may not be running
- Check: `docker ps` - ensure `hangout` container is healthy
- Check hangout logs: `docker logs hangout`

**404 Not Found:**

- Check URL path includes `/rp-api/hangout-service/`
- Verify endpoint exists in backend
- Check NGINX rewrite rules in `conf.d/hangout.conf`

**Cannot connect to NGINX:**

- Ensure Docker Desktop is running
- Check: `docker-compose ps`
- Verify: `curl http://localhost/` returns gateway info

### Modifying Configuration

After changing NGINX config files:

```powershell
# Test configuration is valid
make nginx-test

# Reload without downtime
make nginx-reload

# Or restart NGINX container
docker-compose restart nginx
```

### View Logs

```powershell
# All logs
cd ../../
docker-compose logs -f

# NGINX only
docker logs hangout-nginx -f

# Hangout service only
docker logs hangout -f

# Or use Makefile
cd services/hangout
make nginx-logs
```

---

## 🔐 Phase 2: HTTPS with mkcert (Next Steps)

Coming next after HTTP setup is verified! Will add:

- mkcert local SSL certificates
- HTTPS on port 443
- HTTP/2 support
- Auto HTTP→HTTPS redirect

---

## 📝 Notes

- **Development Mode**: Currently HTTP-only for local testing
- **Gzip**: Removed from Go app, handled by NGINX now
- **Swagger**: Updated to use `/rp-api/hangout-service` base path
- **Docker Network**: Services communicate via internal bridge network
- **Port 9000**: No longer exposed to host, only accessible via NGINX

---

## ✅ Phase 1 Complete Checklist

- [x] NGINX configuration files created
- [x] Docker Compose updated with NGINX service
- [x] Makefile commands added for NGINX management
- [x] Gzip removed from app.go (handled by NGINX)
- [x] Swagger updated with new base path
- [x] HTTP-only reverse proxy working

**Ready to test!** 🚀
