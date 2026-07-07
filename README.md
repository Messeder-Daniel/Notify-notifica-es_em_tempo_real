### Organização atual do backend

Até esta etapa, o backend possui:

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   │   └── health_handler.go
│   └── routes/
│       └── routes.go
├── go.mod
└── go.sum
```

A rota `/health` foi movida para um handler próprio, e o registro das rotas foi centralizado em `internal/routes`.
