# Hotel reservation backend
This project serves as the backend for a hotel reservation system. It's built using MongoDB for the database and GoFiber as the web framework; and is containerized with Docker.
 
- users -> book room from an hotel
- admins -> going to check reservation/bookings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD API -> JSON
- Rooms -> CRUD API -> JSON
- Scripts -> database management -> seeding, migration


## Stack
- MongoDB
- GoFiber
- Docker

to run mongo on docker
```
docker run --name mongodb -d -p 27017:27017 mongo:latest
```

## Project environment variable
```
HTTP_LISTEN_ADDRESS=:3333
JWT_SECRET=
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_URL_TEST=mongodb://localhost:27017
```