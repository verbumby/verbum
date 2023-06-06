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
		--target=es2017 \
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

.PHONY: es-sync
es-sync:
	aws --profile verbum s3 sync s3://verbumby-backup elastic/backup --delete

.PHONY: es-restore-last
es-restore-last:
	elastic/restore-last-backup.bash

.PHONY: es-run
es-run:
	elastic/elasticsearch/bin/elasticsearch \
		-Expack.security.enabled=false \
        -Expack.profiling.enabled=false \
		-Ehttp.host=127.0.0.1 \
		-Ecluster.name=verbum-dev \
		-Enode.name=verbum-1 \
		-Ecluster.initial_master_nodes=verbum-1 \
		-Epath.data=$$(pwd)/elastic/data \
		-Epath.repo=$$(pwd)/elastic/backup

.PHONY: es-setup-backup-repo
es-setup-backup-repo:
	elastic/setup-snapshot.bash $$(pwd)/elastic/backup

