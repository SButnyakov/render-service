FROM golang:1.21.3-bullseye
#ENV CONFIG_PATH=config/
ENV PGDATA=/postgres-storage
ENV CONFIG_PATH=./config/local.yaml
ENV STORAGE_PATH=./storage
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /backend-api
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN go build -C ./cmd/api -o serverApp
WORKDIR /backend-api/cmd/backend-api
CMD ["./serverApp"]
