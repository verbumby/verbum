.PHONY: build
build: fe-build
	go build -v github.com/verbumby/verbum/backend/verbum

.PHONY: build-ci
build-ci:
	go build -v github.com/verbumby/verbum/backend/verbum

.PHONY: run
run: build
	./verbum serve

.PHONY: build-parsers
build-parsers:
	cd backend/dictionary/dslparser && go tool github.com/mna/pigeon -o grammar.go grammar.peg

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
		--platform=node \
		--target=node20 \
		--outdir=frontend/dist
	rm frontend/dist/server.css

	find frontend/dist/public/ -type f \
		\( -name "*.js" -o -name "*.js.map" -o -name "*.css" -o -name "*.css.map" -o -name "*.png" \) \
		-print0 | xargs -0 rm -f

	npx esbuild frontend/browser.tsx frontend/theme.ts \
		--bundle \
		--define:process.env.NODE_ENV='"production"' \
		--minify \
		--sourcemap \
		--platform=browser \
		--target=es2016 \
		--outdir=frontend/dist/public \
		--entry-names=[name]-[hash] \
		--loader:.png=file

	cp -v frontend/index.html frontend/dist/index.html
	cp -v frontend/*.png frontend/dist/public/

.PHONY: es-run
es-run:
	ES_JAVA_OPTS="-Xms1600m -Xmx1600m" \
	elastic/elasticsearch/bin/elasticsearch \
		-Expack.security.enabled=false \
		-Expack.profiling.enabled=false \
		-Ehttp.host=127.0.0.1 \
		-Ecluster.name=verbum-dev \
		-Enode.name=verbum-1 \
		-Ecluster.initial_master_nodes=verbum-1 \
		-Epath.data=$$(pwd)/elastic/data
