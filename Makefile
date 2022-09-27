OUTPUT_DIR=./bin
VERSION=`cat VERSION`
ENDPOINT=http://localhost:8080/v1
GIT_COMMIT=`git rev-list -1 HEAD | cut -c1-8`


raven:
	@CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/raven -ldflags "-X 'main.version=${VERSION}' -X 'main.endpoint=${ENDPOINT}' -X 'main.gitCommit=${GIT_COMMIT}'" ./cmd/raven/main.go

clean:
	-@rm -r $(OUTPUT_DIR)/* 2> /dev/null || true

.PHONY: clean
