DB_URL=postgresql://postgres:cst4Ever@localhost:5432/mypostgres?sslmode=disable

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
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
  .PHONY:  createdb dropdb sqlc  test migrateup migratedown server mock migrateup1 migratedown1 postgres db_docs db_schema proto

