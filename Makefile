.PHONY: build
build:
	go build -v github.com/verbumby/verbum/backend/cmd/verbumsrvr

.PHONY: run
run: build
	./verbumsrvr

.PHONY: build-ctl
build-ctl:
	go build -v github.com/verbumby/verbum/backend/cmd/verbumctl

.PHONY: fe-build-watch
fe-build-watch:
	npx webpack --watch --progress

.PHONY: fe-build-prod
fe-build-prod:
	NODE_ENV=production npx webpack

.PHONY: fe-run
fe-run:
	cd frontend/dist && node server.bundle.js
