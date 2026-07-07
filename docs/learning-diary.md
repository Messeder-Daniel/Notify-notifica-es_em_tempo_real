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

### Perguntas que um professor poderia fazer

#### 1. Por que você separou handlers e routes?

Para separar responsabilidades. A camada de rotas registra os caminhos da API, enquanto os handlers contêm a lógica de resposta das requisições.

#### 2. Qual é a função do `main.go` agora?

Inicializar o router, registrar as rotas e iniciar o servidor.

#### 3. O que é um handler?

É a função ou método que processa uma requisição HTTP e retorna uma resposta.

#### 4. Por que usar a pasta `internal`?

Porque ela indica que aqueles pacotes são internos da aplicação e não devem ser usados como biblioteca externa por outros projetos.

#### 5. O comportamento da rota `/health` mudou?

Não. A resposta continua a mesma. O que mudou foi apenas a organização interna do código.
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

### Perguntas que um professor poderia fazer

#### 1. Para que serve o Docker Compose?

Serve para definir e executar containers necessários ao projeto, como o PostgreSQL.

#### 2. O que é uma variável de ambiente?

É uma configuração externa ao código, usada para adaptar a aplicação a diferentes ambientes.

#### 3. O que é um pool de conexões?

É um conjunto de conexões abertas com o banco que podem ser reutilizadas pela aplicação.

#### 4. O que acontece se o banco estiver fora do ar?

A rota `/health` retorna erro indicando que o banco está desconectado.

#### 5. Por que o PostgreSQL é importante neste projeto?

Porque ele armazenará usuários e notificações de forma persistente.
