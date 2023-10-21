#!/bin/bash
cd ./backend-api
docker build -t backend-cont .
cd ..
cd ./frontend
docker build -t front-cont .
docker run --rm -d -p 3000:3000 -t front-cont
docker run --rm -d --network host backend-cont

