BINARY=whisperchat
.DEFAULT_GOAL := run

run:
	go build && ./$(BINARY)
