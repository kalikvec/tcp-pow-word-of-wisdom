start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go

run:
	docker-compose up --abort-on-container-exit --force-recreate --build