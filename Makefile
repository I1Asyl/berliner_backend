.SILENT:

.DEFAULT_GOAL := run

run: 
	go run .

database_up:
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" up 

database_down: 
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" down 

test_services: 
	cd pkg/services && go test -v -cover

test_repository: 
	cd pkg/repository && go test -v -cover

test_handler: 
	cd pkg/handler && go test -v -cover

