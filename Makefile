HOST = http://localhost:8080

run:
	-make down
	docker compose up --build

down:
	docker compose down --remove-orphans

integration-test:
	@echo "Running integration tests..."
	cd ./cmd && go test -v
	@echo "Integration tests completed."

# Request using 10s token
requestByToken10Sec:
	curl -X GET $(HOST) \
		-H "API_KEY: 2c02b5ce-04d0-4c75-9810-c3e75c397956" \
		-H "Content-Type: application/json"

# Request using 30s token
requestByToken30Sec:
	curl -X GET $(HOST) \
		-H "API_KEY: a6b3fdef-c107-4970-8ecc-94817ed5968c" \
		-H "Content-Type: application/json"

# Request using 2min token
requestByToken2Min:
	curl -X GET $(HOST) \
		-H "API_KEY: 16a661d8-ce97-44b3-a405-a1400d705de8" \
		-H "Content-Type: application/json"

# Request using a regular token
requestByTokenRegular:
	curl -X GET $(HOST) \
		-H "API_KEY: 44acd872-daa3-4d5c-a421-c01119a3d30a" \
		-H "Content-Type: application/json"

# Request using IP
requestByIP:
	curl -X GET $(HOST)
