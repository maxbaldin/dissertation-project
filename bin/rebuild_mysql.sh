cd ..
sudo docker-compose stop -t 1 mysql
sudo docker-compose build mysql
sudo docker-compose up --no-start mysql
sudo docker-compose start mysql