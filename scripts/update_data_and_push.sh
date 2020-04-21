#! /bin/bash

docker build . -f Dockerfile-prod -t covid-data
docker run -v "$HOME/.ssh:/root/.ssh" -it covid-data    
