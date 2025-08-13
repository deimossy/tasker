docker_build:
	docker build -t tasker .

docker_run:
	docker run -p 8080:8080 --name tasker tasker:latest

docker_stop:
	docker stop tasker

docker_remove:
	docker rm tasker

docker_clean:
	docker rmi tasker:latest	

docker_logs:
	docker logs -f tasker

.PHONY: docker_build docker_run docker_stop docker_clean docker_logs