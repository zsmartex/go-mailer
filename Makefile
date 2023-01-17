docker-image:
	docker build -t zsmartex/mailer --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} .
