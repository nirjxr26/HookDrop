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

tekton-apply:
	kubectl wait --for=condition=available --timeout=300s deployment/tekton-pipelines-controller -n tekton-pipelines
	kubectl wait --for=condition=available --timeout=300s deployment/tekton-pipelines-webhook -n tekton-pipelines
	kubectl wait --for=condition=available --timeout=300s deployment/tekton-events-controller -n tekton-pipelines
	kubectl apply -f k8s/tekton/namespace.yaml
	kubectl apply -f k8s/tekton/serviceaccount.yaml
	kubectl apply -f k8s/tekton/tasks.yaml
	kubectl apply -f k8s/tekton/pipeline.yaml

tekton-run:
	kubectl create -f k8s/tekton/pipelinerun.yaml

tekton-delete:
	kubectl delete -f k8s/tekton/pipelinerun.yaml --ignore-not-found
	kubectl delete -f k8s/tekton/pipeline.yaml --ignore-not-found
	kubectl delete -f k8s/tekton/tasks.yaml --ignore-not-found
	kubectl delete -f k8s/tekton/serviceaccount.yaml --ignore-not-found
	kubectl delete -f k8s/tekton/namespace.yaml --ignore-not-found
