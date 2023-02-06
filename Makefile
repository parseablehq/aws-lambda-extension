fmt:
	@go fmt ./{cmd,pkg}/*

vet:
	@go vet ./{cmd,pkg}/*

clean:
	@rm -rf extensions layer*.zip

build: clean fmt vet
	@CGO_ENABLED=0 GOOS="linux" GOARCH=$(or $(GOARCH), "amd64") go build -ldflags "-s -w" -o extensions/parseable-extension cmd/parseable-extension/main.go

compress:
	@upx -9 -q extensions/parseable-extension

package: build
	@zip -9 -q -r layer-$(or $(GOARCH), "x86_64").zip extensions
	@rm -rf extensions

image:
	@docker build \
		-t parseablerepo/parseable-lambda-extension:$(or $(VERSION), 1) \
		-t parseablerepo/parseable-lambda-extension:latest .

publish: package
	@aws lambda publish-layer-version \
		--layer-name parseable-lambda-extension-$(or $(GOARCH), "x86_64")-v1-0 \
		--compatible-runtimes provided provided.al2 nodejs16.x nodejs18.x ruby2.7 java11 java8 go1.x java8.al2 python3.7 python3.8 python3.9 \
        --compatible-architectures $(or $(GOARCH), "x86_64") \
		--description "Lambda function extension for logging to Parseable" \
		--license-info "Apache-2.0" \
		--zip-file fileb://layer-$(or $(GOARCH), "x86_64").zip \
		--output json | tee metadata.json
		