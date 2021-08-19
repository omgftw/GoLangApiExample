# Stock API Example App

## Configuration

Configuration is done via environment variables:

- `BASEURL`: The Base URL for the Stock API to use
- `APIKEY`: The API key to use
- `SYMBOL`: The stock ticker symbol to use
- `NDAYS`: the amount of days to return via a GET request

## Usage

```shell
cd ~/go/src/
git clone https://github.com/omgftw/GoLangApiExample.git
go get
go run main.go
```

## Docker

### Building

```shell
docker build -t stock-api .
```

### Run locally

Docker Hub image:

```shell
docker run --rm -p 8080:8080 omgftw/stock-api
```

Local image:

```shell
docker run --rm -p 8080:8080 stock-api 
```

Open http://127.0.0.1:8080 in your browser

## Kubernetes

```shell
kubectl create secret generic stock-api-secret --from-literal=APIKEY='YOUR_KEY_HERE' 
kubectl apply -f kubernetes/
```

It can be accessed using your kubernetes url at the path /stock-api

Example:
http://192.168.64.3/stock-api
