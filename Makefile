start: 
	@docker-compose up
	
stop:
	@docker-compose rm -v --force --stop
	@docker rmi prod_setup