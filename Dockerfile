FROM golang:1.20 as builder
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o ./mindia

FROM alpine as runner
RUN apk add --no-cache libwebp=1.3.0-r2 libwebp-tools

WORKDIR /app

RUN mkdir -p .bin/webp
RUN ln -s /usr/bin/cwebp .bin/webp/cwebp
RUN ln -s /usr/bin/dwebp .bin/webp/dwebp

RUN addgroup --system --gid 1001 runner
RUN adduser --system --uid 1001 runner

COPY --from=builder --chown=runner:runner /app/mindia .

RUN chmod +x ./mindia

ENV SKIP_DOWNLOAD true

USER runner

EXPOSE 3500

HEALTHCHECK --interval=5s --retries=1 --timeout=500ms CMD curl --fail http://localhost:3500/v0/healthz || exit 1

ENTRYPOINT ["./mindia", "--server"]