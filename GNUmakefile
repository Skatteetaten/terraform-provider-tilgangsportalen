default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... $(TESTARGS) -cover -coverprofile=c.out -timeout 120m
	go tool cover -html=c.out -o coverage.html
