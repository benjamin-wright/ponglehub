.PHONY: cluster repos clean deploy

start: tf-cluster tf-infra trust
stop: untrust tf-infra-clean tf-cluster-rm
restart: stop start

pause: untrust
	k3d cluster stop pongle
	docker stop pongle-registry

resume:
	docker start pongle-registry
	k3d cluster start pongle

tf-cluster:
	cd infra/terraform/cluster && terraform apply -auto-approve

tf-cluster-rm:
	cd infra/terraform/cluster && terraform destroy -auto-approve

tf-infra:
	cd infra/terraform/infra && terraform apply -auto-approve

tf-infra-rm:
	cd infra/terraform/infra && terraform destroy -auto-approve

tf-infra-clean:
	cd infra/terraform/infra && rm -f terraform.tfstate && rm -f terraform.tfstate.backup

trust:
	sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $(shell pwd)/infra/terraform/infra/.scratch/ingress-ca.crt
	npm config set -g cafile $(shell pwd)/infra/terraform/infra/.scratch/ingress-ca.crt
	cp ~/.npmrc ~/.npmrc.bak
	./infra/setup-local.sh
	helm repo add local https://helm.ponglehub.co.uk

untrust:
	sudo security remove-trusted-cert -d $(shell pwd)/infra/terraform/infra/.scratch/ingress-ca.crt || true
	npm config delete -g cafile
	mv ~/.npmrc.bak ~/.npmrc || true
	helm repo remove local || true

geppetto:
	cd tools/geppetto && make install

deploy:
	kubectl get ns | grep ponglehub || kubectl create ns ponglehub
	kubectl annotate namespace ponglehub linkerd.io/inject=enabled --overwrite
	helm dep update deployment
	helm upgrade ponglehub deployment \
		-i \
		--namespace ponglehub