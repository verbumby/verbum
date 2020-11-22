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
	cd frontend/dist && node server.bundle.js --verbum-api-url=https://localhost:8443

.PHONY: es-sync-backup
es-sync:
	aws --profile verbum s3 sync s3://verbumby-backup elastic/backup --delete

.PHONY: es-restore-last
es-restore-last:
	curl -XDELETE 'localhost:9200/access-log-*,dict-*' ; echo
	LAST=$$(curl localhost:9200/_snapshot/backup/_all 2>/dev/null | jq -r '.snapshots[].snapshot' | sort | tail -n 1) \
	&& curl -XPOST localhost:9200/_snapshot/backup/$$LAST/_restore
