# Stock API Example App

## Configuration

Configuration is done via environment variables:

- `BASEURL`: The Base URL for the Stock API to use
- `APIKEY`: The API key to use
- `SYMBOL`: The stock ticker symbol to use
- `NDAYS`: the amount of days to return via a GET request

## Usage

```shell
go get
go run main.go
```

## Docker

### Building

```shell
docker build -t stock-api .
```

### Run locally

```shell
docker run --rm -p 8080:8080 stock-api 
```

Open http://127.0.0.1:8080 in your browser

## Kubernetes

```shell
kubectl create secret generic stock-api-secret --from-literal=APIKEY='YOUR_KEY_HERE' 
kubectl apply -f kubernetes/
```
