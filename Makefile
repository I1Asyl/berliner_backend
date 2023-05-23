.SILENT:

run: 
	go run .

database_up:
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" up 

database_down: 
	migrate -source file://migrations/ -database "mysql://$(USERNAME):$(PASSWORD)@$(PROTOCOL)($(ADDRESS))/$(DBNAME)" down 