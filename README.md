# Membership App

This is a Membership App built with Go. The application allows users to manage memberships via a REST API.

## Prerequisites
You need the following installed to run the application locally:
- Go (version 1.22 or later)
- MongoDB (or Docker for running MongoDB in Docker)


## Running the Application locally

The application stores data in a MongoDB database, so you will need a MongoDB instance running locally. To start a Mongo Docker container, run the following command:

1. **Start a Mongo Docker container**
   ```sh
   docker run -d -p 27017:27017 --name membership-app mongo
   ```

2. **Run the unit tests**
   ```sh
   go test ./...
   ```

3. **Start the Application**
   ```sh
   go run cmd/main.go
   ```


## Running the Application with Docker
As an alternative to running the application locally, you can run the app using Docker. This method ensures that all dependencies and configurations are handled within Docker containers.

1. **Build and start the Application in Docker**

   Use the following command to build the application Docker image and start the services defined in the `docker-compose.yml` file:

   ```sh
   docker-compose up --build
   ```


## Calling the Application's Rest API

Below are some of the sample requests:

### Creating a member
```
curl --location 'localhost:8080/member' \
--header 'Content-Type: application/json' \
--data-raw '{
 "firstName": "Rafael",
 "lastName": "Nadal",
 "email": "Rafael.Nadal@gmail.com",
 "dateOfBirth": "1986-06-03"
}' 
```

### Updating a member by id
```
curl --location --request PUT 'localhost:8080/member/970973' \
--header 'Content-Type: application/json' \
--data-raw '{
 "email": "Rafael.Nadal@tennis.com"
}'
```

### Getting a member by member id
```
curl --location 'localhost:8080/member/970973'
```

### Deleting a member by member id
```
curl --location --request DELETE 'localhost:8080/member/970973'
```