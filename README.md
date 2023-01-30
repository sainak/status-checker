
# Status Checker

Golang http server to check active status of websites


## Run Locally

Clone the project

```shell
git clone https://github.com/sainak/status-checker && cd status-checker
```

Configure environment variables

```shell
cp example.env .env
```

Build the app

```shell
go build
```

Start postgres
```shell
docker run -d --name postgres1 \
    -v postgres_data:/var/lib/postgresql/data \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=website-status \
    -p 5432:5432 \
    postgres:latest
```

Migrate db
```shell
./status-checker migrate up
```

Start the server

```shell
./status-checker serve
```
