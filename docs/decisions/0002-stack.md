# Decision 0002: Use Go Backend and React Frontend

## Status

Accepted

## Decision

Billtap will use:

- Go for backend runtime
- React + TypeScript for frontend
- SQLite as local persistent storage

## Rationale

Go fits:

- HTTP APIs
- webhook workers
- embedded runtime
- Docker packaging
- deterministic local server behavior

React fits:

- hosted checkout UI
- portal UI
- dense developer dashboard
- scenario timeline and debugging interactions

SQLite fits:

- local persistence
- easy reset
- deterministic CI fixtures
- low operational overhead

