#!/bin/bash
set -e

cd $(dirname $0)

docker build --tag accelerator.apps/db:1.0.0 --file ../build/package/db/Dockerfile .
docker container run --publish 5432:5432 --detach --volume ~/.accelerator.apps.postgresql:/var/lib/postgresql/data --env POSTGRES_PASSWORD=password123 --name db accelerator.apps/db:1.0.0

docker build --build-arg service=generator --tag accelerator.apps/generator:1.0.0 --file ../build/package/service/Dockerfile ..
docker container run --publish 8081:9000 --detach --env CALENDAR_TEMPLATE_PATH=/assets/calendar.goics --env GENERATOR_CONFIG_PATH=/assets/generator-config.json --name generator accelerator.apps/generator:1.0.0

sleep 10 # give some time to the database to spin up

docker build --build-arg service=consumer --tag accelerator.apps/consumer:1.0.0 --file ../build/package/service/Dockerfile ..
docker container run --network host --detach --env GENERATOR_SERVICE_URL_FORMAT=http://localhost:8081/events/%s?offset=%d\&limit=%d \
    --env USERS='["jimmy@accelerator-apps.com", "johnny@accelerator-apps.com", "anna@accelerator-apps.com", "maria@accelerator-apps.com", "fabio@accelerator-apps.com"]' \
    --env INTERVAL=60 --env EVENTS_COUNT=5000 --env DB_HOST=localhost --env DB_PORT=5432 --env DB_USER=postgres --env DB_PASSWORD=password123 \
    --name consumer accelerator.apps/consumer:1.0.0
# set INTERVAL=5 if you want to synchronize every 5min
