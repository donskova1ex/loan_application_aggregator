PROJECT?=github.com/donskova1ex/loan_application_aggregator
API_NAME?=app_aggregator
API_VERSION?=0.0.1
API_CONTAINER_NAME?=docker.io/donskova1ex/${API_NAME}


clean_api:
	rm -rf bin/api

api_docker_build:
	docker build --no-cache -t ${API_CONTAINER_NAME}:${API_VERSION} -t ${API_CONTAINER_NAME}:latest -f dockerfile.api .

api_local_run:
	go run ./cmd/api/app_aggregator_api.go

