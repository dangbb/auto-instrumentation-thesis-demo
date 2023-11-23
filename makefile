REPODIR := $(pwd)

generate:
	protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/audit.proto

audit-migrate:
	echo \# make migrate-create name="$(name)"
	go run service/auditservice/cmd/main.go migrate create $(name)

build:
	go build -o bin/audit service/auditservice/cmd/cmd.go
	go build -o bin/customer service/customer/cmd/cmd.go
	go build -o bin/interceptor service/interceptor/cmd/cmd.go
	go build -o bin/warehouse service/warehouse/cmd/cmd.go

run-warehouse:
	KAFKA_BROKER=localhost:9092 KAFKA_TOPIC=warehouse LOGGER_LEVEL=debug MYSQL_HOST=localhost MYSQL_PORT=3320 MYSQL_DBNAME=dbname MYSQL_USERNAME=username MYSQL_PASSWORD=password INTERCEPTOR_ADDRESS=localhost:8090 AUDIT_ADDRESS=localhost:8091 WAREHOUSE_ADDRESS=http://localhost:8092 CUSTOMER_ADDRESS=http://localhost:8093 HTTP_PORT=8092 bin/warehouse server

run-interceptor:
	KAFKA_BROKER=localhost:9092 KAFKA_TOPIC=warehouse LOGGER_LEVEL=debug MYSQL_HOST=localhost MYSQL_PORT=3320 MYSQL_DBNAME=dbname MYSQL_USERNAME=username MYSQL_PASSWORD=password INTERCEPTOR_ADDRESS=localhost:8090 AUDIT_ADDRESS=localhost:8091 WAREHOUSE_ADDRESS=http://localhost:8092 CUSTOMER_ADDRESS=http://localhost:8093 HTTP_PORT=8090 bin/interceptor server

run-customer:
	KAFKA_BROKER=localhost:9092 KAFKA_TOPIC=warehouse LOGGER_LEVEL=debug MYSQL_HOST=localhost MYSQL_PORT=3320 MYSQL_DBNAME=dbname MYSQL_USERNAME=username MYSQL_PASSWORD=password INTERCEPTOR_ADDRESS=localhost:8090 AUDIT_ADDRESS=localhost:8091 WAREHOUSE_ADDRESS=http://localhost:8092 CUSTOMER_ADDRESS=http://localhost:8093 HTTP_PORT=8093 bin/customer server

run-audit:
	KAFKA_BROKER=localhost:9092 KAFKA_TOPIC=warehouse LOGGER_LEVEL=debug MYSQL_HOST=localhost MYSQL_PORT=3320 MYSQL_DBNAME=dbname MYSQL_USERNAME=username MYSQL_PASSWORD=password INTERCEPTOR_ADDRESS=localhost:8090 AUDIT_ADDRESS=localhost:8091 WAREHOUSE_ADDRESS=http://localhost:8092 CUSTOMER_ADDRESS=http://localhost:8093 GRPC_PORT=8091 MIGRATION_FOLDER=service/auditservice/migration bin/audit server

.PHONY: generate migrate-create common-env