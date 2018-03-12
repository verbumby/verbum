FROM golang:1.9 as backend

WORKDIR $GOPATH/src/github.com/verbumby/verbum
COPY vendor vendor/
COPY backend backend/
RUN go install -v github.com/verbumby/verbum/backend/cmd/verbum

FROM debian:stretch
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /verbum
COPY --from=backend /go/bin/verbum /verbum/verbum
COPY statics statics/
COPY templates templates/
CMD [ "/verbum/verbum" ]
