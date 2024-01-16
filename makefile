.PHONY: build run

PORT ?= 8080

build:
	go build -o bank_parser main.go

run:
	PORT=$(PORT) go run main.go