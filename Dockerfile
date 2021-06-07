FROM golang:1.16 as builder
COPY main.go .
ARG SKAFFOLD_GO_GCFLAGS
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /app main.go

FROM alpine:3
ENV GOTRACEBACK=single
CMD ["./app"]
COPY --from=builder /app .
