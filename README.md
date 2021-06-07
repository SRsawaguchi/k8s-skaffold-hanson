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


## k8s resource: pod.yamlの作成

```
kubectl create deployment \
    --image skaffold-example \
    --port 8080 \
    --replicas=1 \
    api \
    --dry-run=client \
    -o yaml > k8s-deployment.yaml
```

※`expose`は対象の`deployment`が存在しない場合`dry-run`でもエラーになる。  
そこで、一度deploymentを`apply`してから以下を実行する。  

```
kubectl apply -f k8s-deployment.yaml

kubectl expose deploy api \
    --port 80 \
    --target-port 8080 \
    --type NodePort \
    --dry-run=client \
    -o yaml > k8s-nodeport.yaml

kubectl delete -f k8s-deployment.yaml
```

## skaffold.yamlの生成

```
skaffold init
```

`? Do you want to write this configuration to skaffold.yaml?`と聞かれるので`yes`で書き込む。  

## build
ここで、`Dockerfile`の内容をもとにImageが作成される。  

```
skaffold build
```

## 開発開始
```
skaffold dev
```

なお、URLは以下のように取得する。  

```
minikube service <service> --url
```

ここでは、

```
minikube service api --url
```

戻ってきたURLに対してcurlなどを行って動作確認を行う。  

```
curl $(minikube service api --url)\?name\=sawa
```

`main.go`を書き換えると自動的にbuildが実行され、k8sクラスタにデプロイされる。  


## exec user process caused: no such file or directoryについて
goで`net/http`を利用していると、`skaffold dev`をしたときにこの不具合が起こる。    

```
 standard_init_linux.go:219: exec user process caused: no such file or directory
 ```

 これはビルド時に環境変数`CGO_ENABLED`を`0`にすれば解決する。  
 netパッケージはデフォルトではcgoが必須であるが、それは通常OSが通信を仲介するから。  
 Alpineでは、それに必要なライブラリが存在しないため、エラーになる。  
 ※単一で動作するのではなくて、netに関しては動的にリンクされるということ。  

 しかし、`CGO_ENABLED=0`の状態だと、pure goのパッケージを利用する。  
 そのため、この問題を回避できる。  

 ```
 The net package requires cgo by default because the host operating system must in general mediate network call setup. On some systems, though, it is possible to use the network without cgo, and useful to do so, for instance to avoid dynamic linking. The new build tag netgo (off by default) allows the construction of a net package in pure Go on those systems where it is possible.
 ```
 https://golang.org/doc/go1.2

