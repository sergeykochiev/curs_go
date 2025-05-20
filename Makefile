run:
	go run main.go
run_reinit:
	rm initialized
	rm main.db
	go run main.go
