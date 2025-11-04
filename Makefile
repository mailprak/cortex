.PHONY: help build build-local docker-build docker-run docker-shell clean test install

# Variables
BINARY_NAME=cortex
DOCKER_IMAGE=cortex
DOCKER_TAG=latest
GO=go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build-local: ## Build cortex binary locally
	$(GO) build -o $(BINARY_NAME) .
	@echo "✓ Built $(BINARY_NAME)"

build: build-local ## Alias for build-local

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✓ Built Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: ## Run cortex in Docker (pass ARGS="your-command")
	docker run --rm -it \
		-v $(HOME)/.kube:/root/.kube:ro \
		-v $(PWD)/example:/cortex/example \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		cortex $(ARGS)

docker-shell: ## Open shell in Docker container
	docker run --rm -it \
		-v $(HOME)/.kube:/root/.kube:ro \
		-v $(PWD)/example:/cortex/example \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		/bin/bash

docker-k8s-example: ## Run K8s health check example in Docker
	docker run --rm -it \
		-v $(HOME)/.kube:/root/.kube:ro \
		-e NAMESPACE=default \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		cortex exec -p /cortex/example/k8s/k8s_cluster_health

docker-system-example: ## Run system health check example in Docker
	docker run --rm -it \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		cortex exec -p /cortex/example/system_health_check

clean: ## Clean built binaries
	rm -f $(BINARY_NAME)
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@echo "✓ Cleaned"

test: ## Run tests
	$(GO) test -v ./...

install: build-local ## Install cortex to /usr/local/bin (requires sudo)
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installed to /usr/local/bin/$(BINARY_NAME)"

# Example neurons and synapses
create-neuron: ## Create a new neuron (NEURON_NAME=name TYPE=check|mutate)
	./$(BINARY_NAME) create-neuron $(NEURON_NAME) -t $(TYPE)

create-synapse: ## Create a new synapse (SYNAPSE_NAME=name)
	./$(BINARY_NAME) create-synapse $(SYNAPSE_NAME)

# Docker Compose targets
compose-up: ## Start services with docker-compose
	docker-compose up -d

compose-down: ## Stop services
	docker-compose down

compose-logs: ## View logs
	docker-compose logs -f

# CI/CD helpers
ci-build: ## Build for CI
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	@echo "✓ Built binaries for multiple platforms"

# Quick examples
example-system: build-local ## Run system health check example
	cd example/system_health_check && ../../$(BINARY_NAME) exec -p .

example-k8s: build-local ## Run K8s health check example
	cd example/k8s/k8s_cluster_health && ../../../$(BINARY_NAME) exec -p .
