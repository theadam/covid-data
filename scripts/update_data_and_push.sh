#! /bin/bash

docker build . -f Dockerfile-prod -t covid-data
docker run -v "$HOME/.ssh:/root/.ssh" covid-data    
docker system prune -a -f
