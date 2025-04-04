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
request-token-10-sec:
	curl -X GET $(HOST) \
		-H "API_KEY: 2c02b5ce-04d0-4c75-9810-c3e75c397956"

# Request using 30s token
request-token-30-sec:
	curl -X GET $(HOST) \
		-H "API_KEY: a6b3fdef-c107-4970-8ecc-94817ed5968c"

# Request using 2min token
request-token-2-min:
	curl -X GET $(HOST) \
		-H "API_KEY: 16a661d8-ce97-44b3-a405-a1400d705de8"

# Request using a regular token
request-token-regular:
	curl -X GET $(HOST) \
		-H "API_KEY: 44acd872-daa3-4d5c-a421-c01119a3d30a"

# Request using IP
request-ip:
	curl -X GET $(HOST)

# Request using random IP
request-ip-random:
	curl -X GET $(HOST) \
		-H "X-Forwarded-For: 192.168.0.1"

# Load testing
## using the default network created by docker compose
load-test-ip:
	docker run --rm --network=go-rate-limiter_default williamyeh/hey -n 15 -c 2 -H "X-Forwarded-For: 192.168.0.1" http://app:8080/

load-test-token:
	docker run --rm --network=go-rate-limiter_default williamyeh/hey -n 15 -c 2 -H "API_KEY: ac76eaf1-4793-430a-a7fa-23716f10ab81" http://app:8080/

load-test-ip-batch:
	@for i in $$(seq 1 10); do \
		IP="192.168.0.$$(shuf -i 1-255 -n 1)"; \
		echo "Sending request $$i from $$IP"; \
		docker run --rm --network=go-rate-limiter_default williamyeh/hey -n 10 -c 5 -H "X-Forwarded-For: $$IP" http://app:8080/; \
		sleep $$(shuf -i 1-3 -n 1); \
	done

load-test-token-batch:
	@for i in $$(seq 1 10); do \
		TOKEN="abc$$(shuf -i 1-100 -n 1)"; \
		echo "Sending request $$i from $$TOKEN"; \
		docker run --rm --network=go-rate-limiter_default williamyeh/hey -n 20 -c 5 -H "API_KEY: $$TOKEN" http://app:8080/; \
		sleep $$(shuf -i 1-3 -n 1); \
	done
