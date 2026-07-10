# Arquitetura do Notify

## Visão geral

Notify é um sistema de notificações em tempo real desenvolvido com Golang, Gin, PostgreSQL, WebSocket, Vite, React, TypeScript e shadcn/ui.

O objetivo do projeto é permitir que usuários autenticados enviem notificações para outros usuários cadastrados, recebam avisos instantâneos no navegador e acompanhem o ciclo de vida dessas notificações.

O sistema possui dois papéis principais:

- `admin`: usuário administrador, com acesso ao painel de gerenciamento de usuários.
- `user`: usuário comum, que pode enviar, receber, responder e concluir notificações.

---

## Stack utilizada

### Backend

- Golang
- Gin
- PostgreSQL
- pgx
- JWT
- bcrypt
- WebSocket

### Frontend

- Vite
- React
- TypeScript
- Tailwind CSS
- shadcn/ui
- Sonner
- Lucide React

### Infraestrutura local

- Docker
- Docker Compose

---

## Estrutura geral do projeto

```text
real-time-notifications/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go
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
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── lib/
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── types.ts
│   ├── package.json
│   └── vite.config.ts
├── docs/
│   ├── architecture.md
│   └── learning-diary.md
├── docker-compose.yml
└── README.md