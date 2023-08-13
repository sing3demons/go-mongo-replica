# go-mongo-replica
docker-compose mongo replica-set

```
make start
go run main.go
```

```
MONGO_URL
```


```generate key
	mkdir -p cert
	openssl genrsa -out cert/id_rsa 4096
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub
	go run genarateKey/gen.go
```