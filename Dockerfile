FROM golang

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY database/migrations /app/database/migrations

RUN go build -o main ./cmd


RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz -O migrate.tar.gz && \
    tar -xvzf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    rm migrate.tar.gz

EXPOSE 8000

# Comando para executar as migrações e iniciar a aplicação
CMD ["sh", "-c", "migrate -path /app/database/migrations -database 'postgres://frete:frete@frete-rapido-database:5432/frete?sslmode=disable' up && ./main"]

