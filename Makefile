DOCKER_USERNAME := brownbarg
APPLICATION_NAME := bumper
APPLICATION_DIR := ./build
APPLICATION_FILEPATH := ./build/bumper
REFVER := $(shell cat .git/HEAD | cut -d "/" -f 3)
SEMVER ?= $(shell cat .git/HEAD | cut -d "/" -f 3)
HELMPACKAGE := "bumper"

.PHONY: build test clean run docker-run docker-build docker-push docker-tag helm-create helm-package helm-test kdeploy

test:
	go test ./... -v

build:
	@rm -rf build
	@echo $(REFVER)
	CGO_ENABLED=0 go build -v --o $(APPLICATION_FILEPATH) -ldflags="-X 'main.Ver=$(REFVER)'" main.go

clean:
	@rm -rf $(APPLICATION_DIR)

run:
	@echo "Running with args: $(ARGS)"
	@rm -rf $(APPLICATION_DIR)
	@go build -o $(APPLICATION_FILEPATH)
	@$(APPLICATION_FILEPATH) $(ARGS)

docker-build:
	@docker build . -t $(APPLICATION_NAME):$(REFVER)

docker-tag:
	@docker tag $(APPLICATION_NAME):$(REFVER) brownbarg/$(APPLICATION_NAME):$(SEMVER)
	@jq '.version = "$(SEMVER)"' config.json > tmp.json && mv tmp.json config.json

docker-run:
	docker run -p 127.0.0.1:8080:8080/tcp $(APPLICATION_NAME):$(REFVER)

docker-push: docker-tag
	@#docker tag $(APPLICATION_NAME):$(REFVER) brownbag/$(APPLICATION_NAME):$(SEMVER)
	@#jq '.version = $(SEMVER)' config.json > tmp.json && mv tmp.json config.json
	@echo $$DOCKER_PASSWORD | docker login -u $(DOCKER_USERNAME) --password-stdin
	@docker push brownbarg/$(APPLICATION_NAME):$(SEMVER)

helm-create:
	@helm create $(HELMPACKAGE)

helm-package:
	@helm package $(HELMPACKAGE)

helm-test:
	@helm lint $(HELMPACKAGE)

kdeploy: docker-build docker-tag docker-push
	@eval $(minikube docker-env)
	@helm install bumper-release $(HELMPACKAGE) --set image.tag="$(SEMVER)"

kremove:
	@helm uninstall bumper-release $(HELMPACKAGE)