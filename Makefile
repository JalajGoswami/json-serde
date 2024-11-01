
start:
	go run cmd/main.go --file=input.json


build:
	go build -o dist/json-serde.bin cmd/main.go


test:
	go test ./pkg/*/
