FROM golang:1.21.9-bullseye
LABEL maintainer="Siddhartha Basu <siddhartha-basu@northwestern.edu>"
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apt-get -qq update \
    && apt-get -yqq install upx
RUN mkdir -p /graphql-server
WORKDIR /graphql-server
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD cmd cmd
ADD internal internal
RUN go build \
    -a \
    -ldflags "-s -w -extldflags '-static'" \
    -installsuffix cgo \
    -tags netgo \
    -o /bin/app \
    cmd/graphql-server/main.go
RUN upx -q -9 /bin/app

FROM gcr.io/distroless/static
COPY --from=0 /bin/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
