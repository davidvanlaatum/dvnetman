
all: api/openapi.yaml web/src/api/index.ts server $(wildcard pkg/mongo/adapt/mock/*)

web/src/api/index.ts: api/openapi.yaml web/src/api/openapi-generator.yaml
	cd web/src/api && openapi-generator generate -i ../../../api/openapi.yaml -g typescript-fetch -c openapi-generator.yaml
ifeq ($(shell uname),Darwin)
	sed -i '' '1s;^;// noinspection all\n// @ts-nocheck\n;' web/src/api/*/*.ts
else
	sed -i '1s;^;// noinspection all\n// @ts-nocheck\n;' web/src/api/*/*.ts
endif
	git add web/src/api

api/openapi.yaml: $(wildcard api/gen/*.go) $(wildcard api/gen/openapi/*.go) $(wildcard api/gen/code/*.go) $(wildcard api/gen/spec/*.go)
	go run ./api/gen
	git add api/openapi.yaml pkg/openapi

server:
	cd web && npm install && npm run build
	go build -o server ./cmd/server

$(wildcard pkg/mongo/adapt/mock/*): $(wildcard pkg/mongo/adapt/*.go)
	go tool mockery

test-go:
	go run ./scripts/no-tests
	go test ./...

test-js:
	cd web && npm run test

test-playwright:
	cd web && npm run e2e-tests

test: test-go test-js test-playwright

lint-js:
	cd web && npm run lint

lint-go:
	golangci-lint run

lint: lint-go lint-js

on-commit: all test lint
	go run ./scripts/extra-files web/src/api/.openapi-generator/FILES web/src/api

add-commit:
	echo 'make on-commit' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
