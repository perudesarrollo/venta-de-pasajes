# Establecer la imagen base
FROM golang:latest

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el go.mod y go.sum al directorio de trabajo
COPY go.mod .
COPY go.sum .

# Descargar las dependencias del módulo global
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Construir la aplicación
RUN go build -o search-service ./cmd/search-service

# Comando predeterminado para ejecutar la aplicación
CMD ["./search-service"]
