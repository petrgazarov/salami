build:
	go build main.go

run:
	./main

build_and_run: build run

init:
	cd iac && terraform init

plan:
	export TF_VAR_openai_api_key=$(OPENAI_API_KEY); \
	export TF_VAR_assumed_role_secret_token=$(shell uuidgen); \
	cd iac && AWS_PROFILE=salami terraform plan

apply:
	export TF_VAR_openai_api_key=$(OPENAI_API_KEY); \
	export TF_VAR_assumed_role_secret_token=$(shell uuidgen); \
	cd iac && AWS_PROFILE=salami terraform apply

destroy:
	cd iac && AWS_PROFILE=salami terraform destroy

validate:
	cd iac && terraform validate