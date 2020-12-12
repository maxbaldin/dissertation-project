cd ..
docker-compose stop -t 1 ui
docker-compose build ui
docker-compose up --no-start ui
docker-compose start ui