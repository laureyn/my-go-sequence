# Étape 1 : build statique de l'application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copie des fichiers de dépendances d'abord (optimise le cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copie du code source
COPY . .

# Compilation de l'app (binaire statique optimisé)
RUN go build -ldflags="-s -w" -o app

# Étape 2 : image minimale finale
FROM alpine:latest

WORKDIR /root/

# Cop
