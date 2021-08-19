FROM golang:1.17-alpine
WORKDIR /go/src/app/

# Get required go module dependencies before adding the binary to prevent unnecessary layer re-creation
COPY go.mod ./
RUN go mod download

COPY main.go ./
RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app . \
