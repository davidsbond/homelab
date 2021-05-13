apply:
	kustomize build $(namespace) -o output.yaml
	kubectl --context=homelab apply -f output.yaml
	rm output.yaml
