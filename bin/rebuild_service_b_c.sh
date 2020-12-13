cd ..
sudo docker-compose stop -t 1 node_2_services_b_c
sudo docker-compose build node_2_services_b_c
sudo docker-compose up --no-start node_2_services_b_c
sudo docker-compose start node_2_services_b_c