antidup: main.go
	go build main.go
	rm antidup
	cp main antidup
	rm main