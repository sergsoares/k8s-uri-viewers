run:
	air

build-and-push:
	 IMAGE_NAME=svc-controller-$(date +%s)
	 docker build -t ttl.sh/svc-controller:1h .
	 docker push ttl.sh/svc-controller:1h

build:
	docker build -t svc-controller .

tag:
	docker tag svc-controller:latest ttl.sh/svc-controller-1717462569:1h
	docker push ttl.sh/svc-controller-1717462569:1h

k3d-import:
	k3d image import svc-controller -c svc-controller

build-import: build k3d-import

apply:
	kubectl apply -f manifests/