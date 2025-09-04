.PHONY: all dev build up down test clean status logs

all: dev

# Inicia todo o ambiente de desenvolvimento
dev:
	@echo "Starting development environment..."
	@docker compose up --build
	@echo "Services started!"

# Constrói as imagens Docker
build:
	@echo "Building images..."
	@docker compose build

# Inicia os serviços Docker
up:
	@echo "Starting services..."
	@docker-compose up -d

# Para os serviços Docker
down:
	@echo "Stopping services..."
	@docker-compose down

# Roda os testes (consulte README.md para mais detalhes)
test:
	@echo "Para executar os testes, com os serviços rodando (make dev):"
	@echo "  cd movies-service/tests/integration && go test -v"
	@echo "  cd ../unit && go test -v"

# Limpa o ambiente Docker (remove contêineres e volumes)
clean:
	@echo "Cleaning environment..."
	@docker-compose down -v
	@docker system prune -f

# Mostra o status dos serviços
status:
	@docker-compose ps

# Mostra os logs dos serviços em tempo real
logs:
	@docker-compose logs -f
