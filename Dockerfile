FROM golang:alpine

WORKDIR /app
COPY . ./

ENV GOOS=linux \
    GO111MODULE=on \
    GOARCH=amd64 \
    CGO_ENABLED=0 \
    GOPROXY=https://goproxy.io,direct
RUN go build -o api.ikurum.cn .

WORKDIR /app/dist
USER root
COPY ../api.ikurum.cn /app/dist/api.ikurum.cn
RUN chmod 777 api.ikurum.cn
RUN mkdir md
EXPOSE 9091
ENTRYPOINT ["/app/dist/api.ikurum.cn"]

# docker buildx build --platform linux/arm64 -t ikurum/api.ikurum.cn .