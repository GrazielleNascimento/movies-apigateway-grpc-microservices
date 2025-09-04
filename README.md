# Movie Microservice API

Sistema completo de microserviços para gerenciamento de filmes usando **Arquitetura Hexagonal**, **Go**, **MongoDB**, **gRPC** e **Docker**.

## 🏗️ Arquitetura

```
┌─────────────────┐    gRPC    ┌─────────────────┐    MongoDB    ┌─────────────────┐
│   API Gateway   │◄----------►│ Movies Service  │◄-------------►│    Database     │
│   (HTTP REST)   │            │     (gRPC)      │               │   (MongoDB)     │
│     :8080       │            │     :50051      │               │     :27017      │
└─────────────────┘            └─────────────────┘               └─────────────────┘
```

### Componentes

- **API Gateway**: Interface HTTP REST que expõe as APIs publicamente
- **Movies Service**: Serviço gRPC que implementa a lógica de negócio
- **MongoDB**: Banco de dados NoSQL para persistência

## 🚀 Começando

## 🚀 Inicialização Rápida

### Pré-requisitos

- Docker e Docker Compose
- Git
- Make (opcional, mas recomendado)

### 1. Clonar o projeto

```bash
git clone https://github.com/GrazielleNascimento/movies-apigateway-grpc-microservices.git
cd movies-apigateway-grpc-microservices
```

### 2. Configurar variáveis de ambiente:
```bash
# Copie o arquivo de exemplo de variáveis de ambiente
cp .env.example .env

# Abra o arquivo .env e ajuste as configurações conforme necessário
# Importante: Substitua your_username e your_password por valores apropriados
# Se estiver usando Docker Compose, você precisará atualizar:
# - localhost -> mongodb (para o serviço MongoDB)
# - localhost -> movies-service (para o serviço gRPC)
```

### 2. Executar com um único comando

```bash
make dev
```

### 3. Executando os Testes

Com os serviços rodando (`make dev`), em outro terminal execute:

```bash
# Testes de integração
cd movies-service/tests/integration
go test -v

# Testes unitários
cd ../unit
go test -v
```

### 4. Verificar se está funcionando

```bash
# Status dos serviços
make status

# Teste rápido
curl http://localhost:8080/health
curl "http://localhost:8080/api/v1/movies?page=1&limit=5"
```

## 📚 Documentação da API

### Endpoints Disponíveis

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/movies` | Lista todos os filmes (paginado) |
| GET | `/api/v1/movies/{id}` | Busca filme por ID |
| POST | `/api/v1/movies` | Cria novo filme |
| DELETE | `/api/v1/movies/{id}` | Remove filme por ID |
| GET | `/health` | Health check |

### Swagger UI

Acesse a documentação interativa em: **http://localhost:8080/swagger/**

### Parâmetros de Query

- **page**: Número da página (padrão: 1)
- **limit**: Itens por página (padrão: 10, máximo: 100)

## 🛠️ Exemplos de Uso via curl

### 1. Listar todos os filmes

```bash
curl -X GET "http://localhost:8080/api/v1/movies?page=1&limit=10" \
  -H "Content-Type: application/json"
```

**Resposta:**
```json
{
  "data": [
    {
      "id": 8,
      "title": "Edison Kinetoscopic Record of a Sneeze (1894)",
      "year": "1894"
    },
    {
      "id": 10,
      "title": "La sortie des usines Lumière (1895)",
      "year": "1895"
    }
  ],
  "total": 34
}
```

### 2. Buscar filme por ID

```bash
curl -X GET "http://localhost:8080/api/v1/movies/8" \
  -H "Content-Type: application/json"
```

**Resposta:**
```json
{
  "data": {
    "id": 8,
    "title": "Edison Kinetoscopic Record of a Sneeze (1894)",
    "year": "1894"
  }
}
```

### 3. Criar novo filme

```bash
curl -X POST "http://localhost:8080/api/v1/movies" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Meu Filme Incrível",
    "year": "2024"
  }'
```

**Resposta:**
```json
{
  "data": {
    "id": 12345,
    "title": "Meu Filme Incrível",
    "year": "2024"
  },
  "message": "movie created successfully"
}
```

### 4. Deletar filme

```bash
curl -X DELETE "http://localhost:8080/api/v1/movies/8" \
  -H "Content-Type: application/json"
```

**Resposta:**
```json
{
  "message": "movie with ID 8 deleted successfully"
}
```

### 5. Health check

```bash
curl -X GET "http://localhost:8080/health"
```

**Resposta:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## 🔧 Comandos do Makefile

| Comando | Descrição |
|---------|-----------|
| `make dev` | Inicia ambiente de desenvolvimento (instala deps + build + up) |
| `make build` | Constrói as imagens Docker |
| `make up` | Inicia todos os serviços |
| `make down` | Para todos os serviços |
| `make test` | Executa testes (api-gateway e movies-service) |
| `make clean` | Remove containers, volumes e limpa ambiente |
| `make status` | Mostra status dos serviços |
| `make logs` | Mostra logs dos serviços em tempo real |

## 🧪 Testes

### Executar todos os testes

```bash
make test
```

### Testes por tipo

```bash
# Testes unitários
cd movies-service && go test -v ./tests/unit/...

# Testes de integração (requer MongoDB)
cd movies-service && go test -v ./tests/integration/...

# Com coverage
cd movies-service && go test -v -race -coverprofile=coverage.out ./...
```

### Executar com Docker

```bash
# Teste em ambiente isolado
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 🏛️ Arquitetura Hexagonal

### Estrutura do Projeto

```
movie-microservice/
├── api-gateway/                    # API Gateway (HTTP REST)
│   ├── cmd/main.go                 # Entry point
│   ├── internal/
│   │   ├── adapters/              # Adapters (HTTP, gRPC)
│   │   │   ├── grpc/client.go     # gRPC client
│   │   │   └── http/handlers/     # HTTP handlers
│   │   ├── core/                  # Business logic
│   │   │   ├── domain/            # Domain entities
│   │   │   ├── ports/             # Interfaces
│   │   │   └── services/          # Business services
│   │   └── config/                # Configuration
│   └── Dockerfile
├── movies-service/                # Movies Service (gRPC)
│   ├── cmd/main.go                # Entry point
│   ├── internal/
│   │   ├── adapters/              # Adapters (gRPC, Database)
│   │   │   ├── grpc/server.go     # gRPC server
│   │   │   └── database/mongodb.go # MongoDB adapter
│   │   ├── core/                  # Business logic
│   │   │   ├── domain/            # Domain entities
│   │   │   ├── ports/             # Interfaces
│   │   │   └── services/          # Business services
│   │   └── config/                # Configuration
│   ├── tests/                     # Tests
│   │   ├── unit/                  # Unit tests
│   │   └── integration/           # Integration tests
│   └── Dockerfile
├── proto/                         # Protocol Buffers
│   └── movies/movies.proto
├── scripts/                       # Initialization scripts
└── docker-compose.yml
```

### Princípios Implementados

1. **Separation of Concerns**: Separação clara entre lógica de negócio e infraestrutura
2. **Dependency Inversion**: Dependências apontam sempre para dentro (interfaces)
3. **Port and Adapters**: Interfaces (ports) e implementações (adapters) bem definidas
4. **Testability**: Facilita criação de testes com mocks
5. **Maintainability**: Código bem organizado e fácil de manter

## 🔄 Comunicação gRPC

### Definição do Serviço

```protobuf
service MovieService {
    rpc GetMovies(GetMoviesRequest) returns (GetMoviesResponse);
    rpc GetMovie(GetMovieRequest) returns (GetMovieResponse);
    rpc CreateMovie(CreateMovieRequest) returns (CreateMovieResponse);
    rpc DeleteMovie(DeleteMovieRequest) returns (DeleteMovieResponse);
}
```

### Testar gRPC diretamente

```bash
# Instalar grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Listar serviços
grpcurl -plaintext localhost:50051 list

# Testar GetMovies
grpcurl -plaintext -d '{"page": 1, "limit": 5}' \
  localhost:50051 movies.MovieService/GetMovies
```

## 🗄️ MongoDB

### Configuração

- **Host**: localhost:27017
- **Database**: movies_db
- **Collection**: movies
- **Autenticação**: admin/password

### Conectar diretamente

```bash
# Via Docker
docker exec -it movies-mongodb mongosh -u admin -p password

# Usar database
use movies_db

# Listar filmes
db.movies.find().pretty()

# Contar documentos
db.movies.countDocuments()
```

### Schema Validation

O MongoDB está configurado com validação de schema:

```javascript
{
  _id: { bsonType: "int", required: true },
  title: { bsonType: "string", required: true },
  year: { bsonType: "string", pattern: "^[0-9]{4}$", required: true }
}
```

## 🚦 Monitoramento

### Health Checks

```bash
# API Gateway
curl http://localhost:8080/health

# Verificar logs
make logs

# Status dos containers
make status
```

### Logs Estruturados

Todos os serviços usam logging estruturado com slog:

```bash
# Ver logs em tempo real
docker-compose logs -f api-gateway
docker-compose logs -f movies-service
docker-compose logs -f mongodb
```

## 🛡️ Tratamento de Erros

### Códigos de Status HTTP

| Status | Descrição | Quando ocorre |
|--------|-----------|---------------|
| 200 | OK | Requisição bem-sucedida |
| 201 | Created | Recurso criado com sucesso |
| 400 | Bad Request | Parâmetros inválidos |
| 404 | Not Found | Recurso não encontrado |
| 500 | Internal Server Error | Erro interno |

### Exemplo de Resposta de Erro

```json
{
  "error": "invalid movie data",
  "message": "year must be between 1800 and current year + 10"
}
```

## 🔧 Desenvolvimento

### Requisitos para Desenvolvimento

```bash
# Instalar Go 1.21+
go version

# Instalar Protocol Buffers
# Ubuntu/Debian
sudo apt-get install protobuf-compiler

# macOS
brew install protobuf

# Instalar plugins Go para protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Configurar ambiente local

```bash
# Clonar repositório
git clone <repo-url>
cd movie-microservice


# Executar testes
make test

```

### Variáveis de Ambiente

#### API Gateway
- `SERVER_PORT`: Porta HTTP (padrão: 8080)
- `MOVIE_SERVICE_GRPC_ADDRESS`: Endereço do Movies Service (padrão: movies-service:50051)
- `READ_TIMEOUT`: Timeout de leitura em segundos (padrão: 10)
- `WRITE_TIMEOUT`: Timeout de escrita em segundos (padrão: 10)

#### Movies Service
- `MONGODB_URI`: String de conexão MongoDB (padrão: mongodb://mongodb:27017)
- `DATABASE_NAME`: Nome do database (padrão: movies_db)
- `GRPC_PORT`: Porta gRPC (padrão: 50051)
- `MAX_POOL_SIZE`: Tamanho máximo do pool MongoDB (padrão: 10)

## 🐛 Troubleshooting

### Problemas Comuns

#### 1. Serviços não inicializam

```bash
# Verificar logs
make logs

# Verificar se as portas estão em uso
netstat -tlnp | grep :8080
netstat -tlnp | grep :50051
netstat -tlnp | grep :27017

# Limpar e reconstruir
make clean build up
```

#### 2. Erro de conexão gRPC

```bash
# Verificar se movies-service está rodando
docker ps | grep movies-service

# Testar conectividade gRPC
grpcurl -plaintext localhost:50051 list
```

#### 3. MongoDB não conecta

```bash
# Verificar MongoDB
docker exec -it movies-mongodb mongosh --eval "db.runCommand('ping')"

# Verificar logs do MongoDB
docker logs movies-mongodb

# Reinicializar MongoDB
docker-compose restart mongodb
```

#### 4. Dados não aparecem

```bash
# Executar inicialização dos dados
make init

# Verificar se dados foram inseridos
docker exec -it movies-mongodb mongosh -u admin -p password --eval "use movies_db; db.movies.countDocuments()"
```

### Logs de Debug

Para habilitar logs mais verbosos:

```bash
# Adicionar ao docker-compose.yml
environment:
  - LOG_LEVEL=debug

# Ou executar com debug
docker-compose up --build
```

## 🚀 Deployment

### Build para Produção

```bash
# Build otimizado
make build

# Tag para produção  
docker tag movie-microservice_api-gateway:latest api-gateway:v1.0.0
docker tag movie-microservice_movies-service:latest movies-service:v1.0.0
```

### Docker Compose Produção

Criar `docker-compose.prod.yml`:

```yaml
version: '3.8'
services:
  api-gateway:
    image: api-gateway:v1.0.0
    environment:
      - SERVER_PORT=8080
      - MOVIE_SERVICE_GRPC_ADDRESS=movies-service:50051
    restart: unless-stopped
    
  movies-service:
    image: movies-service:v1.0.0
    environment:
      - MONGODB_URI=mongodb://user:pass@mongodb:27017/movies_db
      - GRPC_PORT=50051
    restart: unless-stopped
```

## 📝 Licença

Este projeto está licenciado sob a Licença Apache 2.0 - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuição

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📞 Suporte

Para suporte e dúvidas:

- 📧 Email: graziellenascimento454@gmail.com  
- 🐛 Issues: [GitHub Issues](link-para-issues)
- 📖 Wiki: [GitHub Wiki](link-para-wiki)

---

**Desenvolvido com ❤️ usando Go, gRPC, MongoDB e Arquitetura Hexagonal**