FROM node:21-alpine3.17
ENV CONFIG_PATH=/render-service/backend-api/config/
ENV PGDATA=/postgres-storage
ENV CONFIG_PATH=/mnt/c/users/shadow/desktop/git/render-service/backend-api/config/local.yaml
ENV STORAGE_PATH=/mnt/c/users/shadow/desktop/git/render-service/backend-api/
ENV JWT_SECRET_KEY=randomkekv123
WORKDIR /frontend
COPY . .
RUN npm ci
CMD ["npm", "start"]
