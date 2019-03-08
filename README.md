## Local development

Run PostgreSQL instance locally and set the `DATABASE_URL` environment variable.

```sh
# terminal 1
$ docker run -it --rm -e POSTGRES_PASSWORD='' -p 5432:5432 postgres:alpine

# terminal 2
$ export DATABASE_URL='host=localhost port=5432 database=postgres user=postgres sslmode=disable'
$ go run .
```
