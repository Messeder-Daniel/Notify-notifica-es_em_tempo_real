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
### Executar PostgreSQL com Docker Compose

Na raiz do projeto:

```bash
sudo docker compose up -d
```

Verificar se o container está rodando:

```bash
sudo docker compose ps
```

### Configurar variáveis do backend

Copie o arquivo de exemplo:

```bash
cp backend/.env.example backend/.env
```

O arquivo `.env` contém as configurações locais da aplicação e não deve ser enviado ao GitHub.

### Executar backend com conexão ao banco

Entre na pasta do backend:

```bash
cd backend
```

Execute:

```bash
go run ./cmd/api
```

Teste:

```bash
curl http://localhost:8080/health
```

Resposta esperada:

```json
{
  "database": "connected",
  "message": "API is running",
  "status": "ok"
}
```
### Executar migrations

Com o PostgreSQL rodando, execute na raiz do projeto:

```bash
sudo docker exec -i realtime_notifications_postgres psql -U postgres -d notifications_db < backend/migrations/001_create_initial_schema.sql
```

Verificar tabelas criadas:

```bash
sudo docker exec -it realtime_notifications_postgres psql -U postgres -d notifications_db -c "\dt"
```

Verificar usuários iniciais:

```bash
sudo docker exec -it realtime_notifications_postgres psql -U postgres -d notifications_db -c "SELECT name, email FROM users;"
```

Usuários iniciais:

| Nome | Email | Senha |
|---|---|---|
| Alice | `alice@example.com` | `password` |
| Bob | `bob@example.com` | `password` |
