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

# generated writes test Go code to the filesystem
.PHONY: generated
generated: clean
	go test ./wit/bindgen -write

.PHONY: clean
clean:
	rm -rf ./generated/*
