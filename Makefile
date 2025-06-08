.PHONY: install-krr check-brew

install-krr: check-brew
	@echo "Installing krr..."
	brew tap robusta-dev/homebrew-krr
	brew install krr

check-brew:
	@if ! command -v brew &> /dev/null; then \
		echo "Error: Homebrew is not installed. Please install it first."; \
		echo "Visit https://brew.sh for installation instructions."; \
		exit 1; \
	fi

env-up:
	docker compose -f ./test/docker-compose.yml up

env-down:
	docker compose -f test/docker-compose.yml down

push-metrics:
	curl -X POST --data-binary @metrics.prom http://localhost:8428/api/v1/import/prometheus

python-init:
	uv init
	uv venv
	uv pip install pyyaml

python-run:
	python tools/gen-metrics.py

setup-front:
	cd krr-viewer && npm install && npm run dev
