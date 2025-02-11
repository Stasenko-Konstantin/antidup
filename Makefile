front:
	cd web && npx tsc

server:
	go run cmd/server/main.go

run: front server