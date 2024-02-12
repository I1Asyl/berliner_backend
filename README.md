
# Berliner

Backend of the application for sharing your hobbies effectively.

## Prerequisites
- golang 1.20
- MySQL 8.0
## Installation
- Replace `config.yaml` values with your own database information. 
- Create a `.env` file inside configs folder with variables named `DB_PASSWORD` and `JWT_SECRET` for database user password and JWT secret respectively. 
- Run `make database_up` to create all needed tables for the application. 
- Run `make database_down` to delete all existing databases used by application. 
- Run `make get` or `go get` to install all needed libraries. 
- Run `make` to run backend to test it. 
- Run `make build` to build an executable named `bin` inside a folder or use your custom `go build` command.  

