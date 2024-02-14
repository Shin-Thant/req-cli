# Request Cli

Request cli is a small project created using golang. This project is inspired by **`cUrl`**.

The argument list includes:

| Name         | Description                    |
| ------------ | ------------------------------ |
| url          | URL to make request to.        |
| content-type | `Content-Type` request header. |
| method       | Request `Method`.              |
| body         | Request `Body`.                |

## Example

Build an executable file:

```sh
go build -o req-cli
```

Make a request:

```sh
./req-cli -url=https://example.com -method=POST -content-type=application/json -body='{"name":"john doe"}'
```
