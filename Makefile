run: build
	@./bin/pub-sub

build:
	@go build -o "./bin/pub-sub"