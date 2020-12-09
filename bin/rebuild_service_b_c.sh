cd ..
docker-compose stop -t 1 node_2_services_b_c
docker-compose build node_2_services_b_c
docker-compose up --no-start node_2_services_b_c
docker-compose start node_2_services_b_c