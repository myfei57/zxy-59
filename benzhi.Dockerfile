FROM golang:1.23
WORKDIR /app
ENV GOPROXY=off GOSUMDB=off
COPY . .
RUN go build -mod=vendor -o /usr/local/bin/buscharge ./cmd/buscharge
CMD ["buscharge", "-addr", ":8080", "-data", "/app/data", "-web", "/app/web"]
