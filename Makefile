.PHONY: init clean deploy

init:
	@./infra/cluster-start.sh

clean:
	@./infra/cluster-stop.sh

reset:
	@./infra/reset-npm.sh

deploy:
	kubectl get ns | grep ponglehub || kubectl create ns ponglehub
	kubectl annotate namespace ponglehub linkerd.io/inject=enabled --overwrite
	helm dep update deployment
	helm upgrade ponglehub deployment \
		-i \
		--namespace ponglehub