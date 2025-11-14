.PHONY: help install-backend install-frontend install dev-backend dev-frontend dev clean

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@echo "  make install          - Instalar todas las dependencias (backend + frontend)"
	@echo "  make install-backend  - Instalar dependencias del backend"
	@echo "  make install-frontend - Instalar dependencias del frontend"
	@echo "  make dev              - Ejecutar backend y frontend en desarrollo"
	@echo "  make dev-backend      - Ejecutar solo backend en desarrollo"
	@echo "  make dev-frontend     - Ejecutar solo frontend en desarrollo"
	@echo "  make build-backend    - Compilar backend"
	@echo "  make clean            - Limpiar archivos temporales"

install: install-backend install-frontend ## Instalar todas las dependencias

install-backend: ## Instalar dependencias del backend
	@echo "Instalando dependencias del backend..."
	go mod download
	go mod tidy

install-frontend: ## Instalar dependencias del frontend
	@echo "Instalando dependencias del frontend..."
	cd frontend && npm install

dev-backend: ## Ejecutar backend en desarrollo
	@echo "Iniciando servidor backend..."
	go run cmd/server/main.go

dev-frontend: ## Ejecutar frontend en desarrollo
	@echo "Iniciando servidor frontend..."
	cd frontend && npm run dev

build-backend: ## Compilar backend
	@echo "Compilando backend..."
	go build -o bin/litwick cmd/server/main.go

clean: ## Limpiar archivos temporales
	@echo "Limpiando archivos temporales..."
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/
