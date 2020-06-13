.PHONY: build
build:
	go build -v github.com/verbumby/verbum/backend/cmd/verbumsrvr

.PHONY: run
run: build
	./verbumsrvr
