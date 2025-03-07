#!/bin/bash

# Esperar a que Kafka esté listo
echo "Esperando a que Kafka esté disponible..."
cub kafka-ready -b urblog-kafka-1:9092 1 30

# Crear el topic "tweets" si no existe
echo "Creando el topic 'tweets' si no existe..."
kafka-topics --create \
  --topic tweets \
  --bootstrap-server urblog-kafka-1:9092 \
  --partitions 1 \
  --replication-factor 1 \
  --if-not-exists

echo "Topic 'tweets' creado correctamente."

# Mantener el contenedor activo
exec "$@"