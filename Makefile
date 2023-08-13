start:
	docker-compose -f scripts/docker-compose.yml up -d
	docker exec -it mongo1 mongosh --eval "rs.initiate({_id:\"my-replica-set\",members:[{_id:0,host:\"mongo1:27017\"},{_id:1,host:\"mongo2:27018\"},{_id:2,host:\"mongo3:27019\"}]})"
stop:
	docker compose -f scripts/docker-compose.yml down
	rm -rf ./scripts/data
cert:
	mkdir -p cert
	openssl genrsa -out cert/id_rsa 4096
	openssl rsa -in cert/id_rsa -pubout -out cert/id_rsa.pub
	go run genarateKey/gen.go