build:
	docker build --rm -t forum .
	# docker image prune --filter label=stage=builder -f
run:
	docker run --rm --name forum -p 8080:8080 forum
