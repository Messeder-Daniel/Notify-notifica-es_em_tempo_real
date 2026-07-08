# Diário de Aprendizado

Este documento registra os principais aprendizados durante o desenvolvimento do sistema de notificações em tempo real.

## Etapa 1 — Arquitetura inicial e README base

### O que aprendi

Nesta etapa, aprendi que antes de começar a programar é importante planejar a estrutura do projeto.

Também aprendi que o README não é apenas um detalhe visual do GitHub. Ele funciona como a porta de entrada do projeto e deve explicar o que o sistema faz, quais tecnologias usa e como será executado.

Aprendi ainda que uma arquitetura organizada ajuda a separar responsabilidades entre backend, frontend, banco de dados, regras de negócio e comunicação em tempo real.

### Conceitos novos

#### Repositório Git

Um repositório Git registra o histórico do projeto. Ele permite acompanhar mudanças, criar commits e publicar o projeto no GitHub.

#### README

O README é a documentação principal de um projeto. Ele explica o objetivo, tecnologias, instalação, execução, testes e decisões técnicas.

#### Arquitetura de software

Arquitetura de software é a forma como um sistema é organizado. Ela define responsabilidades, separa módulos e facilita manutenção.

#### Backend

Backend é a parte do sistema que roda no servidor. Ele processa regras de negócio, acessa banco de dados e expõe APIs.

#### Frontend

Frontend é a parte visual do sistema, com a qual o usuário interage no navegador.

#### Banco de dados

Banco de dados é onde as informações persistentes ficam armazenadas. Neste projeto, usaremos PostgreSQL.

#### WebSocket

WebSocket é uma tecnologia que permite comunicação contínua entre cliente e servidor. Diferente do HTTP tradicional, ele mantém uma conexão aberta.

### Dificuldades comuns

#### Por que não começar direto pelo código?

Porque sem planejamento o projeto pode ficar desorganizado rapidamente. A documentação inicial ajuda a guiar o desenvolvimento.

#### Por que separar backend e frontend?

Porque cada parte tem responsabilidades diferentes. O backend cuida dos dados e regras. O frontend cuida da interface e interação do usuário.

#### Por que criar uma pasta docs?

Porque alguns documentos são mais detalhados do que o README. O README deve ser direto, enquanto a pasta docs pode guardar explicações maiores.

#### Por que usar commits pequenos?

Porque commits pequenos facilitam entender o histórico do projeto e corrigir problemas caso algo dê errado.


### Resumo da etapa

Nesta etapa, o foco foi preparar a base do projeto. Ainda não houve implementação de código funcional, mas a estrutura documental inicial foi criada para guiar o desenvolvimento. O documento foi construído com auxílio de Inteligência Artifícial. 

## Etapa 3 — Organização inicial do backend

### O que aprendi

Nesta etapa, aprendi a organizar melhor o backend separando responsabilidades.

Antes, o arquivo `main.go` criava o servidor e também definia diretamente a rota `/health`.

Agora, o `main.go` apenas inicia a aplicação e chama a função responsável por registrar as rotas.

A lógica da rota `/health` foi movida para um handler específico.

### Conceitos novos

#### Handler

Um handler é a função ou método responsável por responder uma requisição HTTP.

No projeto, o método `Check` do `HealthHandler` responde a rota `/health`.

#### Rotas

Rotas são os caminhos da API.

Exemplo:

```text
GET /health
```

A camada `routes` centraliza o registro desses caminhos.

#### Separação de responsabilidades

Separar responsabilidades significa evitar que um único arquivo ou função faça coisas demais.

Neste projeto:

- `main.go` inicia o servidor;
- `routes.go` registra as rotas;
- `health_handler.go` responde a requisição de health check.

#### Constructor em Go

Go não tem construtores como algumas linguagens orientadas a objetos.

Mesmo assim, é comum criar funções como:

```go
func NewHealthHandler() *HealthHandler
```

para construir uma struct de forma clara.

### Dificuldades comuns

#### Por que criar uma struct vazia?

Mesmo que `HealthHandler` esteja vazio agora, essa estrutura prepara o projeto para receber dependências no futuro, como services.

#### Por que não deixar tudo no `main.go`?

Porque o projeto vai crescer. Se deixarmos tudo no `main.go`, ele ficará difícil de manter.

#### O que é `internal` em Go?

A pasta `internal` é uma convenção especial do Go. Pacotes dentro de `internal` só podem ser importados por código dentro do mesmo projeto pai. Isso ajuda a proteger partes internas da aplicação.

## Etapa 4 — PostgreSQL com Docker Compose

### O que aprendi

Nesta etapa, aprendi a subir um banco PostgreSQL usando Docker Compose e conectar o backend Go ao banco.

Também aprendi que configurações sensíveis ou variáveis por ambiente devem ficar em arquivos `.env`, enquanto o `.env.example` documenta quais variáveis são necessárias.

A rota `/health`, que antes apenas indicava se a API estava rodando, agora também verifica se o banco de dados está conectado.

### Conceitos novos

#### PostgreSQL

PostgreSQL é um banco de dados relacional usado para armazenar dados persistentes.

Neste projeto, ele será usado para armazenar usuários e notificações.

#### Docker

Docker permite executar aplicações em containers isolados.

Neste projeto, usamos Docker para rodar o PostgreSQL sem precisar instalar e configurar o banco manualmente.

#### Docker Compose

Docker Compose permite definir e executar serviços usando um arquivo `docker-compose.yml`.

Neste projeto, ele define o serviço `postgres`.

#### Variáveis de ambiente

Variáveis de ambiente são configurações externas ao código.

Exemplos:

```text
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/notifications_db?sslmode=disable
```

#### `.env` e `.env.example`

O `.env` contém valores reais usados localmente.

O `.env.example` documenta quais variáveis são necessárias para executar o projeto.

O `.env` não deve ser enviado ao GitHub.

#### Pool de conexões

Pool de conexões é um conjunto de conexões reutilizáveis com o banco de dados.

Isso evita abrir uma nova conexão a cada requisição, melhorando desempenho e controle de recursos.

### Dificuldades comuns

#### Por que usar Docker para o PostgreSQL?

Porque facilita a execução do projeto em qualquer máquina e evita configurações manuais complexas.

#### Por que o Docker precisou de `sudo`?

Porque o usuário atual ainda não tem permissão para acessar o socket do Docker sem privilégios administrativos.

#### Por que não enviar `.env` para o GitHub?

Porque esse arquivo pode conter senhas, tokens e informações sensíveis.

#### Por que testar o banco no `/health`?

Porque a API pode estar rodando, mas o banco pode estar fora do ar. O health check ajuda a identificar esse tipo de problema.

## Etapa 5 — Migrations iniciais

### O que aprendi

Nesta etapa, aprendi a criar a estrutura inicial do banco de dados usando uma migration SQL.

Aprendi que o PostgreSQL organiza dados em tabelas e que o projeto precisa de tabelas para usuários e notificações.

Também aprendi a criar dados iniciais de teste, chamados seeds.

### Conceitos novos

#### Migration

Migration é um arquivo que descreve alterações no banco de dados, como criação de tabelas, índices e dados iniciais.

#### Tabela

Tabela é uma estrutura do banco composta por colunas e linhas.

#### Chave primária

Chave primária é o campo que identifica cada registro de forma única.

#### UUID

UUID é um identificador único usado para evitar colisões e dificultar previsibilidade.

#### Chave estrangeira

Chave estrangeira é um campo que cria relação entre tabelas.

No projeto, `notifications.user_id` aponta para `users.id`.

#### Índice

Índice é uma estrutura que melhora a velocidade de consultas no banco.

#### Seed

Seed é a inserção de dados iniciais para facilitar testes e demonstrações.

### Dificuldades comuns

#### Por que usar migration?

Para que qualquer pessoa consiga recriar a estrutura do banco de forma previsível.

#### Por que usar UUID?

Porque UUIDs são identificadores únicos e menos previsíveis que números sequenciais.

#### Por que criar índices?

Porque o sistema vai consultar notificações por usuário e por status de leitura. Índices ajudam essas consultas a ficarem mais rápidas.

#### Por que usar `ON CONFLICT DO NOTHING`?

Para permitir executar a migration mais de uma vez sem duplicar usuários.


## Etapa 6 — Models do backend

### O que aprendi

Nesta etapa, aprendi a criar models em Go para representar as principais entidades do sistema.

Foram criados models para usuários e notificações.

Também aprendi a diferença entre uma struct usada internamente e uma struct usada como resposta da API.

### Conceitos novos

#### Struct

Struct é uma estrutura de dados em Go que agrupa campos relacionados.

#### Tags JSON

Tags JSON definem como os campos da struct aparecem quando são convertidos para JSON.

#### DTO

DTO significa Data Transfer Object.

É um objeto usado para transportar dados entre camadas ou entre cliente e servidor.

Exemplos criados nesta etapa:

- `LoginRequest`
- `LoginResponse`
- `CreateNotificationRequest`

#### Campo sensível

Campo sensível é qualquer informação que não deve ser exposta publicamente.

No projeto, `PasswordHash` é sensível e por isso usa `json:"-"`.

#### Campo nullable

Campo nullable é um campo que pode não ter valor.

No projeto, `ReadAt` pode ser nulo porque uma notificação pode ainda não ter sido lida.

### Dificuldades comuns

#### Por que criar `UserResponse` se já existe `User`?

Porque `User` contém `PasswordHash`, que não deve ser exposto para o frontend.

#### Por que `ReadAt` é ponteiro?

Porque ele pode ser nulo. Uma notificação não lida ainda não possui data de leitura.

#### Por que usar tags `binding`?

Porque futuramente o Gin poderá validar automaticamente campos obrigatórios das requisições.

## Etapa 7 — Repositories

### O que aprendi

Nesta etapa, aprendi a criar repositories para acessar o PostgreSQL a partir do backend Go.

Aprendi que repositories concentram as queries SQL e evitam que handlers ou services precisem conhecer detalhes do banco de dados.

### Conceitos novos

#### Repository

Repository é uma camada responsável por acessar e manipular dados no banco.

#### Query parametrizada

Query parametrizada usa placeholders como `$1`, `$2` e `$3`.

Isso evita concatenar strings manualmente e reduz risco de SQL Injection.

#### SQL Injection

SQL Injection é uma falha de segurança em que dados maliciosos são inseridos em uma consulta SQL.

Usar parâmetros ajuda a evitar esse problema.

#### QueryRow

`QueryRow` é usado quando esperamos apenas uma linha como resultado.

#### Query

`Query` é usado quando esperamos várias linhas como resultado.

#### Scan

`Scan` copia os valores retornados pelo banco para variáveis ou campos de uma struct.

#### rows.Close

`rows.Close` libera recursos associados ao resultado de uma query.

### Dificuldades comuns

#### Por que o repository retorna erro?

Porque operações de banco podem falhar por vários motivos, como conexão indisponível, SQL errado ou dados inválidos.

#### Por que usar `id::text` nas queries?

Porque no banco o ID é UUID, mas no model Go estamos usando string. O cast para texto simplifica o `Scan`.

#### Por que `MarkAsRead` recebe também o `userID`?

Para garantir que o usuário só consiga alterar notificações que pertencem a ele.

## Etapa 8 — Services

Nesta etapa, criei a camada de services, responsável pelas regras de negócio da aplicação.

Aprendi que repositories acessam o banco, enquanto services aplicam validações e regras antes de chamar os repositories.

Também aprendi sobre bcrypt, usado para comparar senhas com hashes de forma segura.

A aplicação passa a seguir este fluxo:

```text
handlers -> services -> repositories -> database

## Etapa 9 — Autenticação com JWT

Nesta etapa, implementei autenticação com JWT no backend.

Aprendi que o login recebe email e senha, valida a senha com bcrypt e gera um token assinado.

Também aprendi que rotas protegidas usam o header `Authorization` com o formato:

```text
Bearer <token>
Principais arquivos criados ou alterados:

```text
backend/internal/handlers/auth_handler.go
backend/internal/middlewares/auth_middleware.go
backend/internal/services/auth_service.go
backend/internal/routes/routes.go
backend/internal/config/config.go
backend/cmd/api/main.go
```

Testes manuais realizados:

```text
POST /auth/login
GET /auth/me
```

O endpoint `/auth/me` confirmou que o token JWT estava válido e que o middleware conseguiu extrair o ID e o email do usuário autenticado.
