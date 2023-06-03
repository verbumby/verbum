.PHONY: build
build: fe-build
	go build -v github.com/verbumby/verbum/backend/cmd/verbum

.PHONY: build-ci
build-ci:
	go build -v github.com/verbumby/verbum/backend/cmd/verbum

.PHONY: run
run: build
	./verbum serve

.PHONY: build-parsers
build-parsers:
	cd backend/pkg/ctl/dictimport/dictparser/dsl && pigeon -o grammar.go grammar.peg
	cd backend/pkg/dictionary/dslparser && pigeon -o grammar.go grammar.peg

.PHONY: fe-lint
fe-lint:
	npx tsc --noEmit

.PHONY: fe-build
fe-build:
	npx esbuild frontend/server.tsx \
		--bundle \
		--define:process.env.NODE_ENV='"production"' \
		--define:global='globalThis' \
		--minify \
		--sourcemap=inline \
		--format=iife \
		--platform=browser \
		--target=es2020 \
		--outdir=frontend/dist
	rm frontend/dist/server.css

	rm -f frontend/dist/public/*.{js,js.map,css,css.map}
	npx esbuild frontend/browser.tsx \
		--bundle \
		--define:process.env.NODE_ENV='"production"' \
		--minify \
		--sourcemap \
		--platform=browser \
		--target=es2016 \
		--outdir=frontend/dist/public \
		--entry-names=[name]-[hash] \
		--loader:.png=file

	cp frontend/index.html frontend/dist/index.html
	cp frontend/favicon.png frontend/dist/public/favicon.png

.PHONY: es-sync-backup
es-sync:
	aws --profile verbum s3 sync s3://verbumby-backup elastic/backup --delete

.PHONY: es-restore-last
es-restore-last:
	curl -XDELETE 'localhost:9200/dict-*' ; echo
	LAST=$$(curl localhost:9200/_snapshot/backup/_all 2>/dev/null | jq -r '.snapshots[].snapshot' | sort | tail -n 1) \
	&& curl -XPOST localhost:9200/_snapshot/backup/$$LAST/_restore -H 'content-type: application/json' -d '{"indices":"dict-*"}'
