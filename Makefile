dev:
	go run ./...

build:
	go build -o bin/hookdrop ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...

scan:
	trivy fs .

clean:
	rm -rf bin/

docker-build:
	docker build -t hookdrop:local .

docker-run:
	docker compose up --build

docker-stop:
	docker compose down

ci-local:
	act -j build

cluster-up:
	bash scripts/setup-cluster.sh

cluster-down:
	kind delete cluster --name hookdrop

argo-ui:
	kubectl port-forward svc/argocd-server -n argocd 8081:443

argo-pass:
	kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo

policy-apply:
	kubectl apply -f k8s/kyverno/

policy-test:
	kubectl get clusterpolicy
