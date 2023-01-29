https://hub.docker.com/_/redis

[//]: # (docker pull redis)

docker run -e "IP=0.0.0.0" -p 7000-7005:7000-7005 grokzen/redis-cluster:latest

http://redis.io/topics/persistence