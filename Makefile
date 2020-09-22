.PHONY: init clean deploy quick stop

init:
	@./infra/cluster-start.sh all

clean:
	@./infra/cluster-stop.sh

quick:
	@./infra/local-repos/start.sh

stop:
	@./infra/local-repos/stop.sh

deploy:
	kubectl get ns | grep ponglehub || kubectl create ns ponglehub
	kubectl annotate namespace ponglehub linkerd.io/inject=enabled --overwrite
	helm dep update deployment
	helm upgrade ponglehub deployment \
		-i \
		--namespace ponglehub