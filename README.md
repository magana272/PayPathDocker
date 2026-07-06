# PayPath

Personal finance dashboard: track income, expenses, debts, and liquid assets; get computed taxes, net worth, debt payoff timelines, cash flow projections, and AI-powered insights.

**Live demo:** [pay-path-mu.vercel.app](https://pay-path-mu.vercel.app/) — frontend on Vercel, API on Render ([paypath-h36f.onrender.com/api](https://paypath-h36f.onrender.com/api/)), data in MongoDB Atlas.

Monorepo with two apps:

```
backend/    Go JSON REST API (net/http + MongoDB), serves on :8000
frontend/   Next.js 16 + React 19 web app, serves on :3000
docs/       design notes
```

## Quick start

### Docker (whole stack)

```bash
docker compose up --build
```

Brings up MongoDB, the API on http://localhost:8000, and the frontend on http://localhost:3000, seeded with demo data (`user@email.com` / `userpassword`). Optional: export `OPENAI_API_KEY` first to enable the AI endpoints. Mongo data persists in the `mongo_data` volume.

Each app also has its own compose file for running it alone: `backend/docker-compose.yml` (API + MongoDB) and `frontend/docker-compose.yml` (frontend only, expects an API on `NEXT_PUBLIC_API_URL`, default `http://localhost:8000/api`). The stacks all publish the same ports, so run one at a time.

### Manual setup

#### 1. Backend (Go 1.26+)

```bash
cd backend
cp .env.example .env   # or create .env — see variables below
make run               # builds and serves the API on :8000
```

Environment (read from `.env` or the shell):

| Variable | Purpose |
|----------|---------|
| `MONGODB_URI` | MongoDB connection string (Atlas or local) |
| `JWT_SECRET` | Signing key for auth tokens |
| `FRONTEND_URL` | Allowed CORS origin (e.g. `http://localhost:3000`) |
| `HTTP_ADDR` | Listen address (default `:8000`) |
| `PORT` | Alternative to `HTTP_ADDR`; binds `:$PORT` (set automatically by Render) |
| `ENV` | `development` / `production` |
| `OPENAI_API_KEY` | Enables the `/api/ai/*` insight endpoints (optional) |

On first run the backend seeds demo data from `backend/seed/*.csv`. `make reseed` drops the database so it re-seeds on next start.

#### 2. Frontend (Node 24)

```bash
cd frontend
cp .env.example .env.local   # NEXT_PUBLIC_API_URL=http://localhost:8000/api
npm install
make dev                     # next dev on :3000
```

## Architecture

- **Backend** is a layered, feature-folder Go module: thin HTTP handlers call per-feature services (`income`, `expenses`, `debts`, `auth`, `reporting`, `ai/*`, etc.), which depend on repository interfaces over MongoDB, with an in-memory TTL cache and singleflight read-collapsing. JWT + bcrypt auth. See `backend/README.md`.
- **Frontend** is a Next.js App Router app: dashboard, explore (payoff / scenarios / cashflow / AI insights), calendar, and settings pages, calling the API via a small fetch wrapper with client-side caching. See `frontend/README.md`.

## Deployment

| Piece | Where | URL |
|-------|-------|-----|
| Frontend | Vercel | https://pay-path-mu.vercel.app/ |
| Backend API | Render (web service) | https://paypath-h36f.onrender.com/api/ |
| Database | MongoDB Atlas | — |

### Backend on Render

Create a **Web Service** pointed at this repo with root directory `backend`:

- **Build command:** `go build -o paypath ./cmd/api/`
- **Start command:** `./paypath`

Render injects `PORT` at runtime and the server binds to it automatically, so no `HTTP_ADDR` is needed. Set these environment variables in the Render dashboard:

```env
MONGODB_URI=mongodb+srv://...        # MongoDB Atlas connection string
JWT_SECRET=<random secret>
FRONTEND_URL=https://pay-path-mu.vercel.app
ENV=production
OPENAI_API_KEY=sk-...                # optional, enables /api/ai/* endpoints
```

`FRONTEND_URL` is echoed back verbatim as the `Access-Control-Allow-Origin` header — it must be the exact Vercel origin, scheme included, **no trailing slash**, or browsers will reject every API response with a CORS error.

For MongoDB Atlas, allow Render's outbound IPs (or `0.0.0.0/0` for a demo) under Network Access, and use a database user with read/write on the app database.

> **Note:** on Render's free tier the service spins down after ~15 minutes of inactivity, so the first request after idle can take up to a minute while it cold-starts (and re-establishes the Atlas connection).

### Frontend on Vercel

Import the repo in Vercel with root directory `frontend` — Next.js is auto-detected, default build settings work. Set one environment variable:

```env
NEXT_PUBLIC_API_URL=https://paypath-h36f.onrender.com/api
```

Because it's a `NEXT_PUBLIC_*` variable, it's baked in at build time — redeploy after changing it.

## Testing

```bash
cd backend
make test        # go test ./... — Mongo repository tests skip unless MONGODB_URI is set
```