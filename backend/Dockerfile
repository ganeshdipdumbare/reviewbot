FROM golang:1.22.0 AS builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=direct go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o backend .

FROM scratch
COPY --from=builder /build/backend /app/
WORKDIR /app
EXPOSE 5174
CMD ["./backend rest"]