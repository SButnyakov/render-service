FROM golang:1.21.3-bookworm
ENV CONFIG_PATH=/render-service/backend-api/config/
RUN echo $CONFIG_PATH
ENV PGDATA=/postgres-storage
WORKDIR /render-service
COPY . .
RUN go mod -C ./backend-api download && go mod -C ./backend-api verify
RUN go mod -C ./backend-api tidy
RUN go build -C ./backend-api/cmd/backend-api -o $PWD/serverApp
#RUN go build -C ./render-app/ -o $PWD/renderApp
#DOCKER_BUILDKIT=1
RUN apt update
RUN apt-get install -y npm
RUN cd frontend/ && npm ci && cd ..
#RUN apt-get install -y postgresql postgresql-contrib
#RUN systemctl start postgresql.service
#VOLUME /postgres-storage
#RUN psql postgres
RUN ./serverApp &
CMD ["npm", "start", "--prefix", "./frontend/"]
