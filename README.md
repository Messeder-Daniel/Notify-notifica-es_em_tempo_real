# Real-Time Notifications

Sistema de notificações em tempo real usando **WebSockets**, desenvolvido com:

- **Backend:** Go + Gin
- **Frontend:** Vite + TypeScript
- **Banco de dados:** PostgreSQL
- **Tempo real:** WebSocket
- **Autenticação:** JWT
- **Ambiente:** Docker Compose para o PostgreSQL

O projeto permite que um usuário faça login, visualize suas notificações, crie novas notificações, marque notificações como lidas e receba notificações em tempo real pelo navegador.

---

## Funcionalidades

- Autenticação de usuário com JWT.
- API HTTP protegida por token.
- Cadastro/listagem de notificações.
- Marcação de notificações como lidas.
- Persistência em PostgreSQL.
- Conexão WebSocket autenticada.
- Envio de notificações em tempo real.
- Frontend em Vite + TypeScript.
- Dashboard com status da conexão WebSocket.
- Contadores de notificações totais e não lidas.

---

## Arquitetura geral

```text
Frontend Vite + TypeScript
        |
        | HTTP + JWT
        v
Backend Go + Gin
        |
        | SQL
        v
PostgreSQL

Backend Go + Gin
        |
        | WebSocket
        v
Frontend em tempo real
```

Fluxo principal:

```text
1. Usuário faz login no frontend.
2. Backend valida as credenciais.
3. Backend retorna um token JWT.
4. Frontend usa o token para consumir a API.
5. Frontend abre conexão WebSocket com o backend.
6. Usuário cria uma notificação.
7. Backend salva a notificação no PostgreSQL.
8. Backend envia evento WebSocket.
9. Frontend atualiza a lista em tempo real.
```

---

## Estrutura do projeto

```text
real-time-notifications/
├── backend/
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middlewares/
│   │   ├── models/
│   │   ├── repositories/
│   │   ├── routes/
│   │   ├── services/
│   │   └── websocket/
│   ├── migrations/
│   ├── .env.example
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── main.ts
│   │   ├── style.css
│   │   └── types.ts
│   ├── index.html
│   ├── package.json
│   └── tsconfig.json
├── docs/
│   ├── architecture.md
│   └── learning-diary.md
├── docker-compose.yml
├── README.md
└── .gitignore
```

---

## Pré-requisitos

Antes de rodar o projeto, é necessário ter instalado:

- Go
- Node.js
- npm
- Docker
- Docker Compose
- Git

---

## Variáveis de ambiente

O backend usa um arquivo `.env`.

Entre na pasta do backend:

```bash
cd backend
cp .env.example .env
```

Exemplo de configuração:

```env
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/notifications_db?sslmode=disable
JWT_SECRET=change-me-in-production
```

---

## Como rodar o projeto

### 1. Subir o PostgreSQL

Na raiz do projeto:

```bash
docker compose up -d
```

Se seu ambiente exigir permissão de administrador para Docker, use:

```bash
sudo docker compose up -d
```

---

### 2. Executar a migration

Na raiz do projeto:

```bash
docker exec -i realtime_notifications_postgres \
  psql -U postgres -d notifications_db \
  < backend/migrations/001_create_initial_schema.sql
```

Se estiver usando Docker com `sudo`:

```bash
sudo docker exec -i realtime_notifications_postgres \
  psql -U postgres -d notifications_db \
  < backend/migrations/001_create_initial_schema.sql
```

A migration cria as tabelas necessárias e usuários de teste.

---

### 3. Rodar o backend

Em um terminal:

```bash
cd backend
go run ./cmd/api
```

O backend ficará disponível em:

```text
http://localhost:8080
```

Endpoint de verificação:

```text
GET http://localhost:8080/health
```

---

### 4. Rodar o frontend

Em outro terminal:

```bash
cd frontend
npm install
npm run dev
```

O frontend ficará disponível em:

```text
http://localhost:5173
```

---

## Usuário de teste

Use as credenciais abaixo para acessar a aplicação:

```text
E-mail: daniel@example.com
Senha: password
```

---

## Endpoints principais

### Health check

```text
GET /health
```

Verifica se a API está rodando e se o banco está conectado.

---

### Login

```text
POST /auth/login
```

Body:

```json
{
  "email": "daniel@example.com",
  "password": "password"
}
```

Resposta esperada:

```json
{
  "token": "jwt-token"
}
```

---

### Usuário autenticado

```text
GET /auth/me
```

Header:

```text
Authorization: Bearer <token>
```

---

### Listar notificações

```text
GET /notifications
```

Header:

```text
Authorization: Bearer <token>
```

---

### Criar notificação

```text
POST /notifications
```

Header:

```text
Authorization: Bearer <token>
```

Body:

```json
{
  "title": "Nova notificação",
  "message": "Mensagem da notificação"
}
```

---

### Marcar notificação como lida

```text
PATCH /notifications/:id/read
```

Header:

```text
Authorization: Bearer <token>
```

---

### WebSocket

```text
GET /ws?token=<jwt>
```

Evento recebido ao conectar:

```json
{
  "type": "connected",
  "user_id": "...",
  "message": "WebSocket connected"
}
```

Evento recebido quando uma notificação é criada:

```json
{
  "type": "notification.created",
  "data": {
    "id": "...",
    "user_id": "...",
    "title": "...",
    "message": "...",
    "is_read": false,
    "created_at": "..."
  }
}
```

---

## Testes manuais recomendados

Após subir PostgreSQL, backend e frontend:

1. Acesse `http://localhost:5173`.
2. Faça login com o usuário de teste.
3. Verifique se o dashboard exibe `WebSocket conectado`.
4. Crie uma nova notificação.
5. Verifique se ela aparece na lista.
6. Marque a notificação como lida.
7. Verifique se os contadores de total e não lidas são atualizados.
8. Atualize a página e confirme que as notificações continuam salvas.

---

## Validação técnica

### Backend

```bash
cd backend
go test ./...
```

### Frontend

```bash
cd frontend
npm run build
```

---

## Documentação complementar

A documentação do desenvolvimento está disponível em:

```text
docs/architecture.md
docs/learning-diary.md
```

O arquivo `learning-diary.md` registra as etapas de aprendizado, decisões técnicas e explicações dos principais conceitos utilizados.

---

## Tecnologias utilizadas

- Go
- Gin
- PostgreSQL
- pgx
- JWT
- WebSocket
- Docker Compose
- Vite
- TypeScript
- HTML
- CSS

---

## Status do projeto

Projeto funcional e integrado.

Inclui backend, banco de dados, autenticação, WebSocket e frontend consumindo a API em tempo real.
