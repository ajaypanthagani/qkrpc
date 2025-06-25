# QkRPC
An RPC framework based on QUIC protocol for low latency

## Example
There's an example of how to use the library under example directory. Follow the below steps to run the example.

<pre lang="markdown"> bash openssl req -x509 -nodes -newkey rsa:2048 \ -keyout example/keys/key.pem \ -out example/keys/cert.pem \ -days 365 \ -config example/keys/localhost.cnf </pre>

<pre lang="markdown"> go run example/main.go </pre>