#!/bin/bash

# Generate redis.conf using the environment variable (Hardcoded! But! Take it easyyyyyyy!)
sed "s/{{REDIS_PASSWORD}}/$(grep REDIS_PASSWORD .env | cut -d '=' -f2)/" ./redis/redis.conf.template > ./redis/redis.conf

# Start Docker Compose services
docker-compose up
