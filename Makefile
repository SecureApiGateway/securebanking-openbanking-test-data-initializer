service := securebanking-test-data-initializer
repo := europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact
helm_repo := forgerock-helm/secure-api-gateway/securebanking-test-data-initializer/
latesttagversion := latest

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
	$(warning no tag supplied; latest assumed)
	$(eval TAG=latest)
else
	$(eval TAG=$(shell echo $(tag) | tr A-Z a-z))
endif
ifndef setlatest
	$(warning no setlatest true|false supplied; false assumed)
	$(eval setlatest=false)
endif
	env GOOS=linux GOARCH=amd64 go build -o initialize
	@if [ "${setlatest}" = "true" ]; then \
		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${tag} -t ${repo}/securebanking/${service}:${latesttagversion} . ; \
		docker push ${repo}/securebanking/${service} --all-tags; \
    else \
   		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${tag} . ; \
   		docker push ${repo}/securebanking/${service}:${tag}; \
   	fi;

package_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	helm dependency update _infra/helm/${service}
	helm template _infra/helm/${service}
	helm package _infra/helm/${service} --version ${version} --app-version ${version}

publish_helm:
ifndef version
	$(error A version must be supplied, Eg. make helm version=1.0.0)
endif
	jf rt upload  ./*-${version}.tgz ${helm_repo}
