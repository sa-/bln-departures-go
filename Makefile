.PHONY: hafas
hafas:
	go tool oapi-codegen -config hafasClient/config.yaml hafasClient/vbb-hafas-test-api.json