name := securebanking-test-data-initializer
repo := europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact
helm_repo := forgerock-helm/secure-api-gateway/securebanking-test-data-initializer/

.PHONY: all
all: mod build

mod:
	go mod download

build: clean
	go build -o initialize

test:
	go test ./...

test-ci: mod
	$(eval localPath=$(shell pwd))
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
	PATH=$(PATH):${localPath}/pact/bin go test ./...

clean:
	rm -f initialize

docker: clean mod
ifndef tag
	$(warning No tag supplied, latest assumed. supply tag with make docker tag=x.x.x service=...)
	$(eval tag=latest)
endif
	env GOOS=linux GOARCH=amd64 go build -o initialize
	docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${name}:${tag} .
	docker push ${repo}/securebanking/${name}:${tag}

package_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	helm dependency update _infra/helm/${name}
	helm template _infra/helm/${name}
	helm package _infra/helm/${name} --version ${version} --app-version ${version}

publish_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	jf rt upload  ./*-${version}.tgz ${helm_repo}
