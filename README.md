# Sistema de Notificações em Tempo Real

Este projeto é um sistema de notificações em tempo real desenvolvido com Golang, Gin Framework, WebSockets, PostgreSQL, Vite e TypeScript.

O objetivo é permitir que usuários autenticados recebam notificações instantaneamente no frontend por meio de uma conexão WebSocket, mantendo também o histórico das notificações persistido em banco de dados.

## Objetivos do projeto

* Construir uma API backend em Golang usando Gin.
* Implementar autenticação simples com JWT.
* Criar uma conexão WebSocket autenticada.
* Gerenciar clientes conectados em tempo real.
* Enviar notificações instantâneas para usuários conectados.
* Persistir notificações no PostgreSQL.
* Permitir consulta ao histórico de notificações.
* Permitir marcação de notificações como lidas.
* Criar um frontend com Vite e TypeScript.
* Implementar reconexão automática no WebSocket.
* Documentar a arquitetura, execução e decisões técnicas.

## Tecnologias utilizadas

### Backend

* Golang
* Gin Framework
* WebSocket
* PostgreSQL
* JWT
* Docker Compose

### Frontend

* Vite
* TypeScript
* React
* WebSocket API do navegador

### Banco de dados

* PostgreSQL

## Estrutura inicial do projeto

```text
real-time-notifications/
├── backend/
├── frontend/
├── docs/
│   ├── architecture.md
│   └── learning-diary.md
├── README.md
└── .gitignore
```

## Arquitetura planejada do backend

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   ├── handlers/
│   ├── routes/
│   ├── middleware/
│   ├── websocket/
│   └── logger/
├── migrations/
├── go.mod
├── go.sum
└── .env.example
```

## Arquitetura planejada do frontend

```text
frontend/
├── src/
│   ├── components/
│   ├── hooks/
│   ├── pages/
│   ├── services/
│   ├── types/
│   ├── utils/
│   ├── App.tsx
│   └── main.tsx
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .env.example
```

## Funcionalidades planejadas

* Login simples com JWT.
* Listagem de notificações do usuário autenticado.
* Criação de notificações via API HTTP.
* Envio de notificações em tempo real via WebSocket.
* Histórico de notificações.
* Marcação de notificações como lidas.
* Exibição instantânea no frontend.
* Reconexão automática do WebSocket.
* Logs no backend.
* Tratamento de erros.
* Documentação completa.

## Fluxo geral da aplicação

1. O usuário faz login no frontend.
2. O backend valida as credenciais.
3. O backend retorna um token JWT.
4. O frontend usa o token para acessar endpoints protegidos.
5. O frontend abre uma conexão WebSocket autenticada.
6. O backend registra o usuário conectado no Hub WebSocket.
7. Uma notificação é criada via API HTTP.
8. O backend salva a notificação no PostgreSQL.
9. O backend envia a notificação em tempo real se o usuário estiver online.
10. O frontend recebe e exibe a notificação imediatamente.
11. O usuário pode consultar o histórico.
12. O usuário pode marcar notificações como lidas.

## Como executar

As instruções completas serão adicionadas conforme o desenvolvimento avançar.

Ao final do projeto, será possível executar:

```bash
docker compose up
```

E acessar o frontend no navegador.

## Como testar

Os testes manuais e automatizados serão documentados durante o desenvolvimento.

Fluxos mínimos que serão testados:

* login;
* criação de notificação;
* recebimento em tempo real;
* listagem de histórico;
* marcação como lida;
* reconexão WebSocket.

## Endpoints planejados

| Método | Rota                      | Descrição                                 |
| ------ | ------------------------- | ----------------------------------------- |
| GET    | `/health`                 | Verifica se a API está funcionando        |
| POST   | `/auth/login`             | Realiza login                             |
| GET    | `/notifications`          | Lista notificações do usuário autenticado |
| POST   | `/notifications`          | Cria uma nova notificação                 |
| PATCH  | `/notifications/:id/read` | Marca uma notificação como lida           |
| GET    | `/ws`                     | Abre conexão WebSocket autenticada        |

## Decisões arquiteturais

### Separação entre backend e frontend

O backend será responsável por regras de negócio, autenticação, banco de dados e WebSocket.

O frontend será responsável por interface, estado visual, chamadas HTTP e conexão WebSocket no navegador.

### Uso de PostgreSQL

O PostgreSQL será usado para persistir usuários e notificações, garantindo que o histórico continue existindo mesmo após reiniciar o sistema.

### Uso de WebSocket

O WebSocket será usado porque permite comunicação em tempo real entre servidor e cliente, sem que o frontend precise ficar perguntando repetidamente se existem novas notificações.

### Uso de JWT

O JWT será usado para autenticar requisições HTTP e também identificar o usuário durante a conexão WebSocket.

## Status do projeto

* [x] Planejamento inicial
* [x] Arquitetura proposta
* [ ] Backend com Gin
* [ ] Banco PostgreSQL
* [ ] Autenticação JWT
* [ ] API de notificações
* [ ] WebSocket backend
* [ ] Frontend Vite + TypeScript
* [ ] Integração em tempo real
* [ ] README final
* [ ] Revisão final

### Autenticação com JWT

Até esta etapa, foi implementada autenticação com JWT no backend.

Endpoints criados:

    POST /auth/login
    GET /auth/me

O endpoint `POST /auth/login` recebe email e senha, valida as credenciais com bcrypt e retorna um token JWT.

O endpoint `GET /auth/me` é protegido por middleware e retorna os dados básicos do usuário autenticado.

O token deve ser enviado no header:

    Authorization: Bearer <token>

Também foi adicionada a variável de ambiente:

    JWT_SECRET=development-secret-change-me

Em produção, esse valor deve ser trocado por uma chave segura.

### API HTTP de notificações

Até esta etapa, foram implementados endpoints HTTP protegidos por JWT para gerenciamento de notificações.

Endpoints criados:

```text
GET /notifications
POST /notifications
PATCH /notifications/:id/read

### WebSocket backend

Até esta etapa, foi implementado o endpoint WebSocket autenticado:

```text
GET /ws?token=<jwt>