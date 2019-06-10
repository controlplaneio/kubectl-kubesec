V?=0

deploy:
	@mkdir -p ~/.kube/plugins/
	@rm -rf ~/.kube/plugins/kubectl-scan || true
	@go build -o ~/.kube/plugins/kubectl-scan


