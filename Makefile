antidup: main.go
	go build main.go
	cp main antidup
	rm main
	GOOS=windows GOARCH=amd64 go build main.go
