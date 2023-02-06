fmt:
	@go fmt ./{cmd,pkg}/*

vet:
	@go vet ./{cmd,pkg}/*

clean:
	@rm -rf extensions layer*.zip

build: clean fmt vet
	@CGO_ENABLED=0 go build -ldflags "-s -w" -o extensions/parseable-extension cmd/parseable-extension/main.go

compress:
	@upx -9 -q extensions/parseable-extension

package: build compress
	@zip -9 -q -r layer-$(or $(GOARCH), "x86_64").zip extensions
	@rm -rf extensions

image:
	@docker build \
		-t parseablerepo/parseable-lambda-extension:$(or $(VERSION), 1) \
		-t parseablerepo/parseable-lambda-extension:latest .

publish: package
	@aws lambda publish-layer-version \
		--layer-name parseable-extension \
		--compatible-runtimes python3.6 python3.7 python3.8 \
        --compatible-architectures $(or $(GOARCH), "x86_64") \
		--description "Lambda function extension for logging to parseable" \
		--license-info "Apache-2.0" \
		--zip-file fileb://layer-$(or $(GOARCH), "x86_64").zip \
		--output json | tee metadata.json
		