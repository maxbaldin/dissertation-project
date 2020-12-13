sudo docker rm $(docker ps -f status=exited -aq)
sudo docker rmi $(docker images -f "dangling=true" -q)
sudo docker volume rm $(docker volume ls -qf dangling=true)