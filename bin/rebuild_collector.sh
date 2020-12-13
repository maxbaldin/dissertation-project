cd ..
sudo docker-compose stop -t 1 collector
sudo docker-compose build collector
sudo docker-compose up --no-start collector
sudo docker-compose start collector