V?=0

deploy:
	@mkdir -p ~/.kube/plugins/scan
	@go build -o ~/.kube/plugins/scan/scan
	@cp plugin.yaml ~/.kube/plugins/scan/

