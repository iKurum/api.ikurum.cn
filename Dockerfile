FROM --platform=$TARGETPLATFORM golang:alpine

WORKDIR $GOPATH/src/api.ikurum.cn
COPY . .

ENV GOARCH=arm64 \
    GOOS=linux \
    GO111MODULE=on \
    GOPROXY=https://goproxy.io,direct \
    CGO_ENABLED=0

RUN go mod tidy
RUN go build -o api .

EXPOSE 9091
ENTRYPOINT ["./api"]

# docker buildx build -t api --platform=linux/arm64/v8 . --push