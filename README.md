# HTTPStatusFilter

A Go program to filter URLs based on HTTP status codes with optional delay and concurrency control.

## Features

- Filter a list of URLs by specified HTTP status codes.
- Support for specifying multiple status codes or ranges.
- Optional delay between requests to avoid overwhelming servers.
- Configurable concurrency to control the number of parallel requests.

## Installation

Ensure you have Go installed and your `GOPATH` or `GOBIN` set up. Then, run:

```sh
go install github.com/vincentcassiau/HTTPStatusFilter@latest
```

Make sure your GOPATH/bin is in your PATH:
```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage
```sh
HTTPStatusFilter <url_list_file> <status_code(s)> [--delay=<delay_ms>] [--concurrency=<n>]
```
* -f <url_list_file>: A file containing the list of URLs.
* -s <status_code(s)>: One or more status codes or status code ranges.
* -d <delay_ms>: Optional delay in milliseconds between each request.
* -c <n>: Optional number of concurrent requests (default is 5).
* -v: Show version.
* -h: Show help.

## Example
```sh
HTTPStatusFilter -f urls.txt -s 200,301,400-499 -d 500 -c 5
```
This command will filter URLs from urls.txt that return status codes 200, 301, or any status code between 400 and 499, with a 500ms delay between each request and a maximum of 5 concurrent requests.
