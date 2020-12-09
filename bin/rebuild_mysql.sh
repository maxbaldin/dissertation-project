cd ..
docker-compose stop -t 1 mysql
docker-compose build mysql
docker-compose up --no-start mysql
docker-compose start mysql