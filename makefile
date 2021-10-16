create_topics:
	@docker-compose exec kafka-1 kafka-topics --create --topic input --partitions 2 --replication-factor 1 --if-not-exists --zookeeper zookeeper:2181
	
delete_topics:
	@docker-compose exec kafka-1 kafka-topics --delete --topic input --zookeeper zookeeper:2181 --if-exists
	@docker-compose exec kafka-1 kafka-topics --delete --topic goka-input-table --zookeeper zookeeper:2181 --if-exists

docker_container_renew:
	@docker-compose stop
	@docker container prune
	@docker-compose up -d

run:
	@echo "running application"
	@godotenv -f .env go run main.go