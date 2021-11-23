all: peg

test: deps 
	go test -race ./...

deps:
	go mod download;

tools: deps
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

peg:
	cd ent/parser; \
	peg query.peg; \
	goimports -w query.peg.go;

