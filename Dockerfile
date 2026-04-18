FROM golang:1.21-alpine
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1
ENV GOTOOLCHAIN=auto
COPY . .
RUN go build -o server .
EXPOSE 8080
CMD ["./server"]