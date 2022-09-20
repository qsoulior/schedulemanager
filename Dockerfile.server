FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd/server cmd/server
COPY internal/app/config.go internal/app/logger.go internal/app/ 
COPY internal/app/server internal/app/server
COPY internal/controller internal/controller 
COPY internal/entity internal/entity
COPY internal/usecase internal/usecase
COPY internal/repository/mongo.go internal/repository/plan.go internal/repository/
COPY pkg pkg

RUN go build -v -o /app/server /app/cmd/server

FROM alpine:latest

WORKDIR /app


COPY configs configs
COPY --from=build /app/server server

EXPOSE 3000
 
CMD [ "./server", "-c", "configs/dev.json"]