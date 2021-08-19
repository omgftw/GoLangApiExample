FROM golang:1.17-alpine
WORKDIR /go/src/app/

# Get required go module dependencies before adding the binary to prevent unnecessary layer re-creation
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest
WORKDIR /
COPY --from=0 /go/src/app/app /
# The following line can be uncommented to allow testing without a valid API key
# COPY data.json /
ENTRYPOINT ["/app"]
