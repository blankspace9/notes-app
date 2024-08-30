FROM golang:alpine AS builder
WORKDIR /app
COPY . /app
# download packages
RUN go mod download
# building server
RUN go build -tags migrate -o /app/bin/server ./cmd/app
# building migrator
RUN go build -tags migrate -o /app/bin/migrator ./cmd/migrator

FROM alpine:latest
# add certificates for correct work of external api
RUN apk --no-cache add ca-certificates
# copy binaries
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/migrator /app/migrator
# copy config file
COPY config /app/config
# copy migrations files
COPY migrations /app/migrations

# Для использования мигратора можно запускать контейнер с мигратором вручную:
# docker run <image-name> /app/migrator