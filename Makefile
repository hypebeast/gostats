release:
	cd scripts; ./release.sh

deploy:
	cd deployment; ansible-playbook -i hosts_deploy deploy.yml
