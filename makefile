run:
	export $(shell cat .env | xargs) && go run ./cmd/server
