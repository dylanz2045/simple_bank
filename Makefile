DB_URL=postgresql://userdb:cst4Ever@8.134.97.76:5432/ailanzbase

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=cst4Ever -d postgres:14-alpine

createdb:
	docker exec -it postgres createdb --username=postgres --owner=postgres  mypostgres

dropdb:
	docker exec -it postgres dropdb mypostgres

sqlc:
	docker run --rm -v "C:\Users\zdlff\zdl\Project:/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1


server:
	go run main.go

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql  --postgres -o doc/schma.sql  doc/db.dbml

mock:
	mockgen -package mockdb  -destination db/mock/store.go  Project/db/sqlc Store
proto:
	make del
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

install:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc


del: target1 target2 

# 定义第一个目标
target1:
	del /s /q pb\*.go
target2:
	del /s /q doc\swagger\*.swagger.json

.PHONY:  createdb dropdb sqlc  test migrateup migratedown server mock migrateup1 migratedown1 postgres db_docs db_schema proto install del
 
