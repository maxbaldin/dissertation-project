cd ..
sudo docker-compose stop -t 1 node_1_service_a
sudo docker-compose build node_1_service_a
sudo docker-compose up --no-start node_1_service_a
sudo docker-compose start node_1_service_a