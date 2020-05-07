# .PHONY tells Make targets are not associated with files
.PHONY : test fmt tidy

ifeq ($(TRAVIS), true)
  CGO_ENABLED := 0
else
  CGO_ENABLED := 1
endif

test:
	GO111MODULE=on go test --race --cover $$(go list ./...)

fmt:
	@gofmt -l -w $(SRC)

tidy:
	GO111MODULE=on go mod tidy