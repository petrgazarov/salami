build:
	go build main.go

run:
	./main

build_and_run: build run

init:
	cd iac && terraform init

plan:
	cd iac && AWS_PROFILE=salami terraform plan

apply:
	cd iac && AWS_PROFILE=salami terraform apply

destroy:
	cd iac && AWS_PROFILE=salami terraform destroy

validate:
	cd iac && terraform validate