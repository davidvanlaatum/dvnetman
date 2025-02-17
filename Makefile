
all: api/openapi.yaml web/src/api/index.ts server $(wildcard pkg/mongo/adapt/mock/*)

web/src/api/index.ts: api/openapi.yaml web/src/api/openapi-generator.yaml
	cd web/src/api && openapi-generator generate -i ../../../api/openapi.yaml -g typescript-fetch -c openapi-generator.yaml
ifeq ($(shell uname),Darwin)
	sed -i '' '1s;^;// @ts-nocheck\n;' $(wildcard web/src/api/**/*.ts)
else
	sed -i '1s;^;// @ts-nocheck\n;' $(wildcard web/src/api/**/*.ts)
endif


api/openapi.yaml: $(wildcard api/gen/*.go) $(wildcard api/gen/openapi/*.go) $(wildcard api/gen/code/*.go)
	go run ./api/gen

server:
	cd web && npm install && npm run build
	go build -o server ./cmd/server

$(wildcard pkg/mongo/adapt/mock/*): $(wildcard pkg/mongo/adapt/*.go)
	go tool mockery

test:
	go run ./scripts
	go test ./...
