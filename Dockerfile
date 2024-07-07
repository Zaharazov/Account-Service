FROM golang:1.22.5
COPY . /app
WORKDIR /app
RUN go mod download
WORKDIR /app/cmd
RUN go build -o goexec main.go
EXPOSE 8080

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_HOST postgres
ENV POSTGRES_PORT 5432
ENV POSTGRES_DB_NAME postgres

CMD ["./goexec"]