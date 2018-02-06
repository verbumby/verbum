FROM golang:1.9 as backend
RUN wget -O $GOPATH/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 \
    && chmod +x $GOPATH/bin/dep
# TODO: check hash sum
RUN go get -u gopkg.in/reform.v1/...
WORKDIR $GOPATH/src/github.com/verbumby/verbum
COPY Gopkg* ./
RUN dep ensure -v -vendor-only
COPY *.go ./
COPY dict dict/
COPY article article/
COPY tm tm/
RUN (cd dict && go generate) \
    && (cd article && go generate)
RUN go install

FROM node:9 as frontend
WORKDIR /verbum
COPY package*.json ./
RUN npm install
COPY .babelrc webpack.config.js ./
COPY frontend frontend/
RUN npx webpack

FROM debian:stretch
WORKDIR /verbum
COPY --from=backend /go/bin/verbum /verbum/verbum
COPY --from=frontend /verbum/statics/admin.js /verbum/statics/admin.js
COPY templates templates/
CMD [ "/verbum/verbum" ]
