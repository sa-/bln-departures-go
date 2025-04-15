.PHONY: hafas-api-client
hafas-api-client:
	go tool oapi-codegen -config hafasClient/config.yaml hafasClient/vbb-hafas-test-api.json

.PHONY: meteosource-api-client
meteosource-api-client:
	go tool oapi-codegen -config meteoSource/config.yaml meteoSource/openapi.json

.PHONY: run
run:
	go run *.go