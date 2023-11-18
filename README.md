<h1 align="center">
  Octgopus
</h1>
<p align="center">
  Trading platform API for submitting and matching orders. 🐙
</p>
<p align="center">
  Built with Golang and the Echo framework. Consumed by a <a href="https://github.com/richo225/orderbook" target="_blank">React and Typescript app</a>.
</p>

<p align="center">
  <a href="https://octgopus.up.railway.app" target="_blank">
    <img src="https://img.shields.io/website?label=backend&&up_message=live&down_message=down&url=https%3A%2F%2Foctgopus.up.railway.app%2F" />
  </a>
  <a href="https://octgopus.up.railway.app" target="_blank">
    <img src="https://img.shields.io/website?label=frontend&&up_message=live&down_message=down&url=https%3A%2F%2Foctgopus.up.railway.app%2F" />
  </a>
  <a href="https://github.com/richo225/octgopus/actions/workflows/ci.yml" target="_blank">
    <img src="https://github.com/richo225/octgopus/actions/workflows/ci.yml/badge.svg" />
  </a>
  <a href="https://github.com/richo225/octgopus/blob/master/LICENSE.txt" target="_blank">
    <img src="https://img.shields.io/badge/license-MIT-blue.svg" />
  </a>
</p>

![screenshot](https://github.com/richo225/orderbook/blob/main/assets/main-view.png?raw=true)

## Frontend

:art: [richo225/orderbook](https://github.com/richo225/orderbook/) :art:

## Demo

:star: [octgopus.up.railway.app](https://octgopus.up.railway.app/) :star:

## Installation

To run the server, follow these steps:

1. Ensure that Go is installed (version 1.16+ is recommended).
2. Clone the repository:
```shell
  git clone https://github.com/richo225/octgopus.git
```
3. Navigate to the project directory:
```shell
   cd octgopus
```

### Setup

Copy the example .env.dist file in the project directory to .env and fill in the values:

```shell
  cp .env.dist .env
```

```shell
  PORT=<Port the server should run at> eg. 8080
  ALLOWED_ORIGINS=<Request source of the react app for CORS protection> eg. http://localhost:3000
```

### Build
To build the application, run the following command in the project directory:

```
  go build
```
This will create an executable file in the same directory.

### Run
Either run the generated executable via:

```shell
  ./octgopus
```

Or combine the previous steps with:

```
  go run .
```

### Tests
To run the tests, use the following command:

```shell
  go test ./...
```

This will run all tests in the project.