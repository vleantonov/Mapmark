include .env
export


migrate:
	go run cmd/migrator/main.go

run:
	go run cmd/mapmark/main.go