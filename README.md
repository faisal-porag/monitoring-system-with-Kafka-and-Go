### Real Time Monitoring Service Demo Using Golang & Kafka

- Run producer by following command
```shell
make run_producer
```

- Run consumer by following command 
```shell
make run_consumer
```

- CURL
```shell
curl --location 'http://localhost:8083/send-data' \
--header 'Content-Type: application/json' \
--data '{
    "message": "Cash back offer 10% in recharge"
}'
```

- Run kafka using docker compose 
```shell
docker-compose up -d
```

- Check the running container
```shell
docker ps
```

