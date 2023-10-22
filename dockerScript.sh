#!/bin/bash
cd ./backend-api
docker build -t backend-cont .
cd ..
cd ./frontend
docker build -t front-cont .
docker run -d -p 3000:3000 -t front-cont:latest
docker run -d -p 8080:8080 -t backend-cont:latest

