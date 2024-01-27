# Movie Suggestion Service

This project uses Go, Kafka and Cassandra to receive watched movies from a Rest API, send it to Kafka and consume it later to create suggestions and send it to Cassandra where it can be retrieved through a rest api later.

- API at localhost:8080
- Kafka UI at localhost:8081

## How to Start
Run 
```
docker-compose up
```
Feed database with movies from IMDB
```
docker exec -it cassandra1 cqlsh -f schema/cassandra/0_init.cql
```
## APIs

- `GET /movies?page={pageState}` to fetch movies from database
- `GET /users/{id}/suggestions` to fetch suggestions based om watched movies
- `POST /users/{id}/movies/{movieId}/watched` save movie as watched to trigger movie suggestion

## TO DO
- Swagger
- K6 load test