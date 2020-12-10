apply:
	kustomize build $(namespace) -o output.yaml
	kubectl --context=homelab --namespace=$(namespace) apply -f output.yaml
	rm output.yaml
