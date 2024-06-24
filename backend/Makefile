start: 
	docker-compose -f docker-compose.yml up --build -d
stop:
	docker-compose down
test:
	go test ./... -cover
generate:
	go generate ./...
generate-docs:
	swag init 