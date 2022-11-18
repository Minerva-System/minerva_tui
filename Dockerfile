FROM golang:1.19.1-alpine3.16 AS builder
RUN mkdir /minerva-tui
WORKDIR /minerva-tui
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY . .
RUN go build

FROM alpine:3.16 as certs
RUN apk --no-cache add ca-certificates


FROM scratch
ENV TERM=xterm-256color
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /minerva-tui/minerva_tui minerva_tui
ENTRYPOINT ["./minerva_tui"]

