wit_files = $(sort $(shell find testdata -name '*.wit' ! -name '*.golden.*'))

# json recompiles the JSON intermediate representation test files.
.PHONY: json
json: $(wit_files)

.PHONY: $(wit_files)
$(wit_files):
	wasm-tools component wit -j --all-features $@ > $@.json

# golden recompiles the .golden.wit test files.
.PHONY: golden
golden: json
	go test ./wit -update

# generate writes generated Go files in tests and generated directories
.PHONY: generate
generate: clean json
	go test ./wit/bindgen -write
	go generate ./tests

.PHONY: clean
clean:
	rm -rf ./generated/*

# test runs Go and TinyGo tests
GOTESTARGS :=
.PHONY: test
test:
	go test $(GOTESTARGS) ./...
	GOARCH=wasm GOOS=wasip1 go test $(GOTESTARGS) ./...
	tinygo test $(GOTESTARGS) ./...
	tinygo test -target=wasip1 $(GOTESTARGS) ./...
	tinygo test -target=wasip2 $(GOTESTARGS) ./...
	tinygo test -target=wasip2 $(GOTESTARGS) ./tests/...
