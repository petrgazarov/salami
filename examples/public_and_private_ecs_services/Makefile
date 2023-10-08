run:
	../../salami compile

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
	export TF_VAR_openai_api_key=$(OPENAI_API_KEY); \
	export TF_VAR_assumed_role_secret_token=$(shell uuidgen); \
	cd iac && AWS_PROFILE=salami terraform destroy

validate:
	cd iac && terraform validate