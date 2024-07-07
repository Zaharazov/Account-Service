# Базовый образ, над которым производим доработки
FROM golang:1.19
# Определяем рабочую директорию
WORKDIR /app
# Устанавливаем зависиомости
COPY go.mod go.sum ./ 
# здесь копируем requirements
RUN go mod download 
# а здесь их устанавливаем
# Копируем код исходников
COPY *.go ./cmd ./internal .
# Собираем
RUN CGO_ENABLED=0 GOOS=windows go build -o /goexec
# Пробрасываем внешний порт для подключения извне
EXPOSE 8080
# Запускаем
CMD ["/goexec"]