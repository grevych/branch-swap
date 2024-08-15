PATH := $(PATH):$(PWD)/bin


.PHONY: build test integration-test setup-repo clean


build:
	mkdir -p ./bin
	go build -o ./bin/brnchswppr ./cmd/branchswapper.go

setup-repo:
	chmod u+x ./scripts/setup-repo.sh
	./scripts/setup-repo.sh

test:
	go test -v -tags brnchswppr_test ./...

integration-test: build setup-repo
	export PATH
	chmod u+x ./scripts/test.sh
	./scripts/test.sh

clean:
	rm -rf ./tests || true
