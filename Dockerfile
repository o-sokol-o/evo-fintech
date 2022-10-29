FROM golang:1.18 AS builder

RUN go version

COPY . /github.com/o-sokol-o/evo-fintech/
WORKDIR /github.com/o-sokol-o/evo-fintech/

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

FROM alpine:latest

RUN apt-get update

# Add docker-compose-wait tool -------------------
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

WORKDIR /root/

COPY --from=builder /github.com/o-sokol-o/evo-fintech/.bin/app .

CMD ["./app"]