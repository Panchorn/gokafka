# gokafka

# Run local
### Start Kafka & Zookeeper
$ docker-compose -f docker/kafka/docker-compose.yaml up -d

### Start MySQL
$ docker-compose -f docker/mysql/docker-compose.yaml up -d

### Start Redis
$ docker-compose -f docker/redis/docker-compose.yaml up -d

### Download Dependencies & Start Application Manually
$ cd apiA
$ go mod download
$ go run .

$ cd ../streamA
$ go mod download
$ go run .

$ cd ../streamB
$ go mod download
$ go run .

# CURL
## transfers
$ curl --location 'http://localhost:8000/transfers' \
--header 'Content-Type: application/json' \
--data '{
"refId": "refNo104",
"fromId": "3",
"toId": "2",
"amount": 70,
"secretToken": "asdsa"
}'

## get transactions
$ curl --location --request GET 'http://localhost:8000/transfers/transactions' \
--header 'Content-Type: application/json' \
--data '{
"refId": "refNo025",
"fromID": "3",
"toID": "2",
"amount": 70,
"secretToken": ""
}'

# Build image
$ cd apiA
$ make build

$ cd ../streamA
$ make build

$ cd ../streamB
$ make build

$ cd ..
$ docker-compose up -d
