cd ..
sudo docker-compose stop -t 1 ui
sudo docker-compose build ui
sudo docker-compose up --no-start ui
sudo docker-compose start ui