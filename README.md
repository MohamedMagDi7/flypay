# FlyPay Challenge

This repository contains a Go program that represents an API that read the payment transaction data from files (Json format) and return them in the API response as json format.
There are two payment providers `flypayA` and `flypayB`.

Using these 2 files as a source of data for the this task
- `flypayA` data is stored in [flypayA.json](./models/data/flypayA.json)
- `flypayB` data is stored in [flypayB.json](./models/data/flypayB.json)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You should have Golang installed at your device.

you can install from [here](https://golang.org/doc/install)

### Installing using go
After Installing Golang go to the project root directory

```
$ cd go/src/flypay
```

And install dependency libraries

```
$ go get ./...
```

then build the project

```
$ go build
```

then start the server 

```
$ go-docker
```
then open a browser tab with this url [http://localhost:8000/api/payment/transaction](http://localhost:8000/api/payment/transaction)
you should get a list of data provided by providerA and providerB

### Building and Running the Docker image
As we have `Dockerfile` defined in the root directory

```
$ cd go/src/flypay
```

Building the image

```
$ docker build -t go-docker .
```

You can list all the available images by typing the following command -

```
$ docker image ls
```

Running the Docker image 

```
$ docker run -d -p 8000:8000 go-docker
```

Finding Running containers 

```
$ docker container ls
```

Interacting with the app running inside the container

```
$ curl http://localhost:8000/api/payment/transaction
```

## Running the tests

Use this command to run unit tests on the project
 
Go to you root directory 
```
$ go test
```

### Benchmark

you can run benchmarks using this command

```
$ go test -bench=.
```

