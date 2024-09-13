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
	cd tests && go generate ./

.PHONY: clean
clean:
	rm -rf ./generated/*
