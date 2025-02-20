front:
	cd web && spago bundle-app

server:
	go run cmd/server/main.go

run: front server
