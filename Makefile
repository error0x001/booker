APP_NAME := booker
MANIFEST_DIR=.k8s/manifests

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(APP_NAME) ./cmd/booker

run: build
	./bin/$(APP_NAME)

test:
	go test -race -v ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint: lint-install
	golangci-lint run --config=.golangci.yml

deploy:
	kubectl apply -f $(MANIFEST_DIR)

delete:
	kubectl delete -f $(MANIFEST_DIR) --ignore-not-found

redeploy: delete deploy
