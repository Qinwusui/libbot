FROM golang:1.20.2 as builder

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 \
    && GOOS=linux \
    && GOARCH=amd64 \
    && go env -w GOPROXY=https://goproxy.cn,direct &>/dev/null \
    && go build -o libbot &>/dev/null

FROM alpine:latest
COPY --from=builder /build/libbot /libbot/
WORKDIR /libbot
ENV PATH /libbot:$PATH

RUN chmod 777 libbot
CMD ./libbot
