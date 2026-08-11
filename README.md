# Notify — Sistema de Notificações em Tempo Real

Sistema full stack de notificações em tempo real, desenvolvido como projeto prático de backend e frontend.

O Notify permite que usuários autenticados enviem e recebam notificações em tempo real, acompanhem o status das mensagens, respondam notificações e consultem históricos.

## 🚀 Principais funcionalidades

- Cadastro e autenticação de usuários
- Autenticação com JWT
- Controle de acesso por perfil (`admin` e `user`)
- Gerenciamento de usuários
- Envio de notificações
- Recebimento em tempo real via WebSocket
- Notificações visuais com Sonner
- Histórico de notificações enviadas
- Marcação como lida/não lida
- Marcação como concluída/reabertura
- Respostas às notificações
- Interface responsiva

## 🏗️ Arquitetura

```text
Frontend
React + TypeScript + Vite
        │
        │ HTTP / WebSocket
        ▼
Backend
Go + Gin
        │
        ▼
PostgreSQL

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

- React
- TypeScript
- Vite
- Tailwind CSS
- shadcn/ui
- Sonner
- Lucide React

### Infraestrutura

- Docker
- Docker Compose

---

Usuários de demonstração

Para facilitar os testes locais, o banco cria usuários de demonstração na primeira inicialização.

Papel	E-mail	Senha
Admin	admin@example.com	Demo@2026
User	user@example.com	Demo@2026

> Essas credenciais são destinadas exclusivamente ao ambiente de demonstração local.

---

▶️ Como executar

### Pré-requisitos

Go
Node.js
npm
Docker
Docker Compose

1. Clone o projeto
git clone https://github.com/Messeder-Daniel/Notify-notifica-es_em_tempo_real.git
cd Notify-notifica-es_em_tempo_real

2. Inicie o PostgreSQL
docker compose up -d

3. Execute o backend
cd backend
go run ./cmd/api

4. Instale as dependências do frontend

Em outro terminal:

cd frontend
npm install

5. Execute o frontend
npm run dev

Depois, acesse a aplicação pelo endereço exibido pelo Vite.

🧠 O que este projeto demonstra

O Notify foi desenvolvido para praticar conceitos importantes de desenvolvimento de software, incluindo:

### APIs REST

autenticação e autorização
comunicação em tempo real com WebSocket
persistência com PostgreSQL
integração entre frontend e backend
organização de código
desenvolvimento com Git
execução de ambientes locais com Docker

📚 Aprendizados

Este projeto permitiu aprofundar conhecimentos em desenvolvimento backend com Go, comunicação em tempo real, banco de dados relacional e integração com uma aplicação frontend em React e TypeScript.

🔮 Próximos passos

- Evoluir testes automatizados
- Melhorar observabilidade
- Adicionar novas regras de notificação
- Avaliar deploy em ambiente cloud


👨‍💻 Desenvolvido por Daniel Messeder

