FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o barn .

FROM alpine as barn
WORKDIR /
COPY --from=builder /app/barn .

EXPOSE 8888
EXPOSE 32227/udp


ENTRYPOINT ["/barn"]