.PHONY: cluster repos clean deploy

geppetto:
	cd tools/geppetto && make install

cluster:
	@./infra/cluster/start.sh
	geppetto watch

repos:
	@./infra/repos/start.sh
	geppetto watch

clean:
	@./infra/cluster/stop.sh
	@./infra/repos/stop.sh

deploy:
	kubectl get ns | grep ponglehub || kubectl create ns ponglehub
	kubectl annotate namespace ponglehub linkerd.io/inject=enabled --overwrite
	helm dep update deployment
	helm upgrade ponglehub deployment \
		-i \
		--namespace ponglehub