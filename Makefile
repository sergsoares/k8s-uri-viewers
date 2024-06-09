.PHONY: help

# Show this help.
help:
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t

# Use Golang air to live reload with golang
run:
	air

build-and-push:
	 IMAGE_NAME=k8s-uri-viewers-$(date +%s)
	 docker build -t ttl.sh/k8s-uri-viewers:1h .
	 docker push ttl.sh/k8s-uri-viewers:1h

build:
	docker build -t k8s-uri-viewers .

tag:
	docker tag k8s-uri-viewers:latest ttl.sh/k8s-uri-viewers-1717462569:1h
	docker push ttl.sh/k8s-uri-viewers-1717462569:1h

# used to import files
build-import: build k3d-import

# apply manifest folder
apply:
	kubectl apply -f manifests/

# k3d action to create cluster for tests
k3d-create-cluster:
	k3d cluster create k8s-uri-viewers --api-port 127.0.0.1:6500 -p "127.0.0.1:8081:80@loadbalancer"

# k3d action to delete test cluster
k3d-delete-cluster:
	k3d cluster delete k8s-uri-viewers

# k3d import to put image inside k3d cluster
k3d-import:
	k3d image import k8s-uri-viewers -c k8s-uri-viewers

dockerhub-publish:
	docker tag k8s-uri-viewers sergsoares/k8s-uri-viewers:0.0.1
	docker push sergsoares/k8s-uri-viewers:0.0.1

dockerhub-login:
	docker login  -u sergsoares
