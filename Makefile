all: build

build:
	docker build -t michaelpeterswa/honeypot-ingestion .
publish:
	docker push michaelpeterswa/honeypot-ingestion