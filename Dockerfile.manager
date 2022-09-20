FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/manager cmd/manager
COPY internal/app/config.go internal/app/logger.go internal/app/ 
COPY internal/app/manager internal/app/manager
COPY internal/entity internal/entity
COPY internal/usecase internal/usecase
COPY internal/repository internal/repository
COPY pkg pkg

RUN go build -v -o /app/manager /app/cmd/manager

FROM alpine:latest

WORKDIR /app

COPY configs configs
COPY --from=build /app/manager manager

CMD [ "./manager", "-c", "configs/dev.json"]