
include services/kubernetes/etc/help.mk
include services/kubernetes/etc/cli.mk

.PHONY: all
all: services configurations ##@setup deploy everything

.PHONY: deploy
deploy: services configurations ##@setup apply all applications and configurations

.PHONY: services
services: ##@setup apply all applications
	cd services/service-repo && make deploy
	cd services/service-repo2 && make deploy
	cd services/service-repo3 && make deploy
	cd services/service-repo4 && make deploy
	cd services/service-repo5 && make deploy

.PHONY: configurations
configurations: ##@setup apply all configurations
	$(CLI) kubectl apply -R \
		-f configurations/service-repo \
		-f configurations/service-repo3 \
		-f configurations/service-repo5 \

	cd configurations/service-repo2 && make deploy
	cd configurations/service-repo4 && make deploy
