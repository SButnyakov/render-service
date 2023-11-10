FROM golang:1.21.3-bullseye
#ENV CONFIG_PATH=config/
ENV PGDATA=/postgres-storage
ENV API_CONFIG_PATH=/backend-api/config/localapi.yaml
ENV AUTH_CONFIG_PATH=/backend-api/config/localauth.yaml
ENV API_STORAGE_PATH=/backend-api/internal/storage/
ENV AUTH_STORAGE_PATH=/backend-api/internal/storage/
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /backend-api
COPY . .
RUN go mod download && go mod verify
RUN go mod tidy
RUN go build -C ./cmd/api -o ./../serverAPI
RUN go build -C ./cmd/api -o ./../serverAuth
WORKDIR /backend-api/cmd/
CMD ["./serverAPI", "&&", "./serverAuth"]
