cd ..
docker-compose stop -t 1 node_1_service_a
docker-compose build node_1_service_a
docker-compose up --no-start node_1_service_a
docker-compose start node_1_service_a