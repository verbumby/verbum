.PHONY: build
build:
	go build -v github.com/verbumby/verbum/backend/cmd/verbumsrvr

.PHONY: run
run: build
	./verbumsrvr

.PHONY: build-ctl
build-ctl:
	go build -v github.com/verbumby/verbum/backend/cmd/verbumctl
