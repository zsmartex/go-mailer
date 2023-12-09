.PHONY: build
build:
	docker build -t zsmartex/mailer --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} .

.PHONY: push
push:
	docker push zsmartex/mailer