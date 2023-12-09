#!/bin/bash

./redis-server &
./serverAPI &
./serverAuth 
./serverBuff
