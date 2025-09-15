# Terraloom Core API
This is a core api for my miniature e-commerce project Terraloom. The core API serves as the backend of the e-commerce webapp and handles :
- Account registration
- Transaction process
- Payment process

## Specification

The core api written in GO. With specification : 
- Go (ver 1.24.2)
- Gin for web server (ver 1.10.1)
- GORM for object relational mapping and database access (ver 1.30.1)

## Prepare the Database
Using PostgreSQL is reccomended, but any other RDBMS is fine. You can find the ddl for every table in this repository.

## Run in Local
To run this project in local or anyother machine : 
- Make sure GO & GIT installed
```
go --version 
```
- Pull this repository
```
git pull {thisrepositoryurl}
```
- Make sure that your database already running
- Create a new .env file in the root of the project repository
```
PORT=
ENV=
DATABASE_HOST=
DATABASE_PORT=
DATABASE_USER=
DATABASE_PASS=
DATABASE_NAME=
JWT_SECRET=
```
- Fill up the database based on your setup
- The **PORT** part is where the service going to run, make sure the port is free
- Fill the **JWT_SECRET** with your own secret
- To setup Gin server in **release mode** fill the **ENV** with **PRODUCTION** , to setup it in **debug mode** fill the **ENV** with **LOCAL** or **DEV**
- Below is the example of .env
```
PORT=8080
ENV=PRODUCTION
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=terraloom
DATABASE_PASS=1234pass
DATABASE_NAME=terraloom
JWT_SECRET=verysecuresecretnooneknows
```
- Finally to run the service
```
go run ./cmd/api/main.go
```