.PHONY: deps test wire install
.SILENT: wire build clean test
SHELL := bash

build_dir := ./.build
package_dir := ${build_dir}
env := ${app_env}

pkg_path="${GOPATH}/pkg/mod/github.com"
goseanto_repo='nuvemex/goseanto'
goseanto_version=$(shell cat go.mod | grep -o '/goseanto v[0-9].[0-9].[0-9]' | cut -d' ' -f2)
goseanto_resources="${goseanto_version}/resources/config"

define with-env
	@bash -c 'set -o allexport; source .env; set +o allexport; $(1)'
endef

clean:
	rm -rf ${build_dir}

deps:
	go get -u ./...

build: clean
	if [ "${env}" = "" ]; then echo "Please set app_env"; exit 1; fi;
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${build_dir}/search/bootstrap ./lambda/search/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${build_dir}/hinter/bootstrap ./lambda/hinter/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${build_dir}/details/bootstrap ./lambda/details/main.go

	mkdir -p ${build_dir}/resources/config

	# copy goseanto configs
	cp ${pkg_path}/${goseanto_repo}@${goseanto_version}/resources/config/config.yml ${build_dir}/resources/config/goseanto.yml
	cp ${pkg_path}/${goseanto_repo}@${goseanto_version}/resources/config/${env}.yml ${build_dir}/resources/config/goseanto-${env}.yml

	cp ./resources/config/config.yml ${build_dir}/resources/config/
	cp ./resources/config/${env}.yml ${build_dir}/resources/config/${env}.yml

	# duplicate resources in every lambda dir
	cp -R ${build_dir}/resources ${build_dir}/search/resources
	cp -R ${build_dir}/resources ${build_dir}/hinter/resources
	cp -R ${build_dir}/resources ${build_dir}/details/resources

package: build
	@cd ${build_dir}/search && zip -q -r search.zip ./bootstrap ./resources/ && mv search.zip ../
	@cd ${build_dir}/hinter && zip -q -r hinter.zip ./bootstrap ./resources/ && mv hinter.zip ../
	@cd ${build_dir}/details && zip -q -r details.zip ./bootstrap ./resources/ && mv details.zip ../

test:
	go mod download
	cp ${pkg_path}/${goseanto_repo}@${goseanto_version}/resources/config/config.yml ./resources/config/goseanto.yml
	cp ${pkg_path}/${goseanto_repo}@${goseanto_version}/resources/config/testing.yml ./resources/config/goseanto-testing.yml

	@app_env=testing AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_REGION=eu-central-1 \
	go test -v \
			-coverprofile .testCoverage.txt -count=1

	chmod 0777 resources/config/goseanto*

wire:
	wire .

install:
	@./resources/install.sh

install-resources:
	@./resources/install-resources.sh

elasticsearch:
	$(call with-env,go run ./cli/elasticsearch.go)

docker-build:
	@docker-compose build

docker-up:
	@docker-compose up -d

docker-stop:
	@docker-compose stop

docker-logs:
	@docker-compose logs -f
