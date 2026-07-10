# Notify — Sistema de Notificações em Tempo Real

Notify é um sistema de notificações em tempo real desenvolvido com **Golang**, **Gin**, **PostgreSQL**, **WebSocket**, **Vite**, **React**, **TypeScript** e **shadcn/ui**.

O projeto permite que usuários autenticados enviem notificações para outros usuários, recebam avisos em tempo real no navegador, respondam notificações, acompanhem status de leitura, conclusão e histórico de mensagens enviadas.

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

## Funcionalidades

- Cadastro de usuários
- Login com JWT
- Atualização de perfil
- Alteração de senha
- Controle de papéis: `admin` e `user`
- Painel administrativo para gerenciar usuários
- Envio de notificações para usuários cadastrados
- Recebimento de notificações em tempo real via WebSocket
- Toasts em tempo real com Sonner
- Listagem de notificações recebidas
- Histórico de notificações enviadas
- Marcar notificação como lida ou não lida
- Marcar notificação como concluída ou reabrir
- Responder notificações
- Visualização de remetente e destinatário
- Registro do momento de leitura e conclusão
- Interface responsiva

---

## Usuários iniciais

Ao subir o banco do zero, o projeto cria automaticamente dois usuários para teste:

| Papel | E-mail | Senha |
|---|---|---|
| Admin | messederdaniel@outlook.com | Teste@2026 |
| User | barretodaniel11971@hotmail.com | Teste@2026 |

O usuário admin pode acessar a tela de gerenciamento de usuários e alterar papéis entre `admin` e `user`.

---

## Pré-requisitos

Antes de rodar o projeto, instale:

- Go
- Node.js
- npm
- Docker
- Docker Compose

---

## Como rodar o projeto

### 1. Clonar o repositório

```bash
git clone https://github.com/Messeder-Daniel/Notify-notifica-es_em_tempo_real.git
cd Notify-notifica-es_em_tempo_real
