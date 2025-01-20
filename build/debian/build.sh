#!/bin/bash

UGID="$(id -u):$(id -g)" docker-compose up --build
UGID="$(id -u):$(id -g)" docker-compose down
