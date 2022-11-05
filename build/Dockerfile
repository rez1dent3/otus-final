FROM golang:1.19 as build

ENV BIN_FILE "/usr/local/bin/app"
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}
COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/imgproxy/*

FROM alpine:3.9

# healthcheck
RUN apk update
RUN apk add curl

ENV BIN_FILE "/usr/local/bin/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/imgproxy/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE}
