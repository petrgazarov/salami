run:
	../../salami compile

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