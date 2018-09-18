
include ../kubernetes/etc/help.mk
include ../kubernetes/etc/cli.mk

deploy: applications configurations ##@setup apply all applications and configurations

applications: ##@setup apply all applications
	cd ../service-repo && make deploy
	cd ../service-repo2 && make deploy
	cd ../service-repo3 && make deploy
	cd ../service-repo4 && make deploy
	cd ../service-repo5 && make deploy

configurations: ##@setup apply all configurations
	$(CLI) kubectl apply -R \
		-f service-repo \
		-f service-repo3 \
		-f service-repo5 \

	cd service-repo2 && make deploy
	cd service-repo4 && make deploy
