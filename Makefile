.PHONY: cluster repos clean deploy

geppetto:
	cd tools/geppetto && make install

cluster: rustcc
	@./infra/cluster/start.sh

repos: rustcc
	@./infra/repos/start.sh

clean:
	@./infra/cluster/stop.sh
	@./infra/repos/stop.sh

rustcc:
	docker build -t rustcc tools/rust-cross-compiler

deploy:
	kubectl get ns | grep ponglehub || kubectl create ns ponglehub
	kubectl annotate namespace ponglehub linkerd.io/inject=enabled --overwrite
	helm dep update deployment
	helm upgrade ponglehub deployment \
		-i \
		--namespace ponglehub