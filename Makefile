.SILENT:

.DEFAULT_GOAL := run

get: 
	go get

build: 
	go build -o bin . 

run: 
	go run .

database_up:
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" up 

database_down: 
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" down 

test_services: 
	cd pkg/services && test -v -cover




