cd ..
docker-compose stop -t 1 collector
docker-compose build collector
docker-compose up --no-start collector
docker-compose start collector