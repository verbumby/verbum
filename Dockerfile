FROM golang:1.11 as backend

WORKDIR $GOPATH/src/github.com/verbumby/verbum
COPY vendor vendor/
COPY backend backend/
RUN go install -v github.com/verbumby/verbum/backend/cmd/verbumsrvr
RUN go install -v github.com/verbumby/verbum/backend/cmd/verbumctl

FROM debian:stretch
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /verbum
COPY --from=backend /go/bin/verbumsrvr /verbum/verbumsrvr
COPY --from=backend /go/bin/verbumctl  /verbum/verbumctl
COPY statics statics/
COPY templates templates/
CMD [ "/verbum/verbumsrvr" ]
