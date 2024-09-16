FROM golang:1.23 AS builder
WORKDIR /build
COPY . .

# Отключаем CGO и собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app ./cmd/avito_tender_api/main.go

FROM scratch
WORKDIR /bin

ARG POSTGRES_CONN
ARG SERVER_ADDRESS

ENV POSTGRES_CONN=$POSTGRES_CONN
ENV SERVER_ADDRESS=$SERVER_ADDRESS

ENV POSTGRES_CONN=$POSTGRES_CONN
ENV SERVER_ADDRESS=$SERVER_ADDRESS


COPY --from=builder /build/app .
EXPOSE 8080

ENTRYPOINT [ "app" ]
