FROM golang:1.22.3-alpine3.20 AS build
WORKDIR /app
COPY . .
RUN go build -ldflags "-s -w" -o /app/main .

FROM alpine:3.20
COPY --from=build /app/main /app/main
COPY static/ ./static
# Set environment variables here
EXPOSE 8080
ENV GIN_MODE release
CMD ["/app/main"]