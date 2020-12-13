cd ..
sudo docker-compose stop -t 1 node_3_service_d
sudo docker-compose build node_3_service_d
sudo docker-compose up --no-start node_3_service_d
sudo docker-compose start node_3_service_d