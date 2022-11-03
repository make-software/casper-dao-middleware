#!/bin/sh

migrate -database "mysql://$DATABASE_URI" -path /app/resources/migrations up
