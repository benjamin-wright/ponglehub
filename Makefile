.PHONY: cluster repos clean deploy

start: tf-cluster tf-infra trust
stop: untrust tf-infra-clean tf-cluster-rm
restart: stop start

pause: untrust
	k3d cluster stop pongle
	docker stop $(shell docker ps -q --filter name=pongle-) || true

resume:
	docker start $(shell docker ps -aq --filter name=pongle-) || true
	k3d cluster start pongle
	sleep 3
	make trust

tf-init:
	cd infra/terraform/registries && terraform init
	cd infra/terraform/cluster && terraform init
	cd infra/terraform/infra && terraform init

tf-repos:
	cd infra/terraform/registries && terraform apply -auto-approve

tf-repos-rm:
	cd infra/terraform/registries && terraform destroy -auto-approve

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

config:
	./infra/setup-npm.sh
	helm repo add local http://localhost:5002

trust: config
	sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $(shell pwd)/infra/terraform/infra/.scratch/ingress-ca.crt

untrust:
	sudo security remove-trusted-cert -d $(shell pwd)/infra/terraform/infra/.scratch/ingress-ca.crt || true
	./infra/restore-npm.sh
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
