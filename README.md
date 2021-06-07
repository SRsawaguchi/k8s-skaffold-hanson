# k8s-skaffold-hanson

[Skaffold](https://skaffold.dev/)というツールを使ってみる。  
これを使えばローカル環境における開発でもk8sを利用しやすくなる。  

## Goアプリを準備
適当なGoアプリを作成する。  
このようなアプリにする。  

```
> curl http://localhost:8080
Hello, World

> curl http://localhost:8080?name=tanaka
Hello, tanaka!
```

## Dockerfileを作成する

```dockerfile
FROM golang:1.16 as builder
COPY main.go .
ARG SKAFFOLD_GO_GCFLAGS
RUN go build -gcflags="$SKAFFOLD_GO_GCFLAGS" -o /app main.go

FROM alpine:3
ENV GOTRACEBACK=single
CMD ["./app"]
COPY --from=builder /app .
```

[マルチステージビルド](https://matsuand.github.io/docs.docker.jp.onthefly/develop/develop-images/multistage-build/)を利用する。  
先頭の`golang:1.16 as builder`はこのGoアプリをコンパイルする専用のコンテナ。なので、`builder`という名前をつけている。  
その下の`alpine:3`は実際にアプリを稼働させるコンテナ。  
`builder`でコンパイルして生成したバイナリ`/app`をコピーしてくるのが、`COPY --from=builder /app .`。  

これを使うと、複数のファイルでコンテナを生成するスクリプトを書かなくてもいいし、Dockerfileも1つにまとめられる。  

