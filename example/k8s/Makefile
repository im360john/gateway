NAMESPACE=demo
GATEWAY_RELEASE_NAME=gateway

.PHONY: prepare-postgres install-postgres upgrade-postgres uninstall-postgres get-password install-gateway upgrade-gateway uninstall-gateway

install-postgres:
	kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm install postgres bitnami/postgresql -f ./values.postgres.yaml -n $(NAMESPACE)

upgrade-postgres:
	helm upgrade postgres bitnami/postgresql -f ./values.postgres.yaml -n $(NAMESPACE)

uninstall-postgres:
	helm uninstall postgres -n $(NAMESPACE)
	kubectl delete pvc -l app.kubernetes.io/instance=postgres -n $(NAMESPACE)

get-password:
	echo "PostgreSQL Password:"
	kubectl get secret --namespace $(NAMESPACE) postgres -o jsonpath="{.data.postgresql-password}" | base64 --decode && echo

install-gateway:
	helm install $(GATEWAY_RELEASE_NAME) ../../helm/gateway -f ./values.gateway.yaml -n $(NAMESPACE)

upgrade-gateway:
	helm upgrade $(GATEWAY_RELEASE_NAME) ../../helm/gateway -f ./values.gateway.yaml -n $(NAMESPACE)

uninstall-gateway:
	helm uninstall $(GATEWAY_RELEASE_NAME) -n $(NAMESPACE)

template-gateway:
	helm template $(GATEWAY_RELEASE_NAME) ../../helm/gateway -f ./values.gateway.yaml -n $(NAMESPACE) > gateway_template.yaml
