# Etapa de construcción
FROM golang:1.24-alpine as builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

WORKDIR /app

# Copiar los archivos go.mod y go.sum y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto de los archivos del proyecto
COPY . .

# Compilar la aplicación
RUN go build -o urblog ./cmd/

# Crear una imagen más pequeña para ejecutar la aplicación
FROM gcr.io/distroless/base-debian10

WORKDIR /app

COPY --from=builder /app/urblog /app/urblog

# Establecer el comando de inicio
CMD ["/app/urblog"]