
include ../kubernetes/etc/help.mk
include ../kubernetes/etc/cli.mk

make deploy: #@setup apply all configurations
	kubectl apply -r -f service-repo service-repo3 service-repo5
	$(MAKE) -C service-repo2 deploy
	$(MAKE) -C service-repo4 deploy

