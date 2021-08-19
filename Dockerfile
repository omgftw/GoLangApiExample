FROM golang:1.17-alpine
COPY go.mod ./
RUN go get -d -v
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app . \
