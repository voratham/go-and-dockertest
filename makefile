test-int:
	ginkgo -v ./...

clear-container:
	docker rm -f $(docker ps -a -q)

watch-docker:
	watch docker ps

gen-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out