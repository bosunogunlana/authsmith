run:
	export $(shell cat .env | xargs) && go run ./cmd/server

seed:
	export $(shell cat .env | xargs) && go run ./cmd/seed
