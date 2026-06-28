RELEASE := kolesa
NAMESPACE := kolesa
VALUES := -f values.yaml -f values-minikube.yaml

deploy: ## Первый деплой чарта в minikube
	helm install $(RELEASE) . -n $(NAMESPACE) --create-namespace $(VALUES)

redeploy: ## Передеплой (helm upgrade) после правок в чарте/сервисах
	helm upgrade $(RELEASE) . -n $(NAMESPACE) $(VALUES)
