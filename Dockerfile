FROM golang:1.19.1-alpine3.16 AS builder
RUN mkdir /minerva-tui
WORKDIR /minerva-tui
COPY . .
RUN go build

FROM alpine:3.16
ARG APP_DIR=/usr/src/app
ENV TZ=Etc/UTC
ENV APP_USER=appuser
ENV TERM=xterm-256color
RUN addgroup -g 1000 $APP_USER \
    && mkdir -p $APP_DIR \
    && adduser -u 1000 -G $APP_USER -h $APP_DIR -D $APP_USER \
    && chown -R $APP_USER:$APP_USER $APP_DIR
WORKDIR $APP_DIR
COPY --from=builder /minerva-tui/minerva_tui minerva_tui
USER $APP_USER
CMD ["./minerva_tui"]

