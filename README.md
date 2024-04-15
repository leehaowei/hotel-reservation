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