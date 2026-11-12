FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod ./
COPY internal ./internal
COPY main.go ./

RUN go test ./... \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/personal-finance-tracker .

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/personal-finance-tracker /app/personal-finance-tracker
EXPOSE 8080
ENTRYPOINT ["/app/personal-finance-tracker"]
