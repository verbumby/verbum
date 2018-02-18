FROM golang:1.9 as backend

WORKDIR $GOPATH/src/github.com/verbumby/verbum
COPY vendor vendor/
COPY backend backend/
RUN go install -v github.com/verbumby/verbum/backend/cmd/verbum

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
