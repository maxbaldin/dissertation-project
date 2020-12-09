cd ..
docker-compose stop -t 1 node_3_service_d
docker-compose build node_3_service_d
docker-compose up --no-start node_3_service_d
docker-compose start node_3_service_d