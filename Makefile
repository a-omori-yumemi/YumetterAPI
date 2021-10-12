swag.gen: swagger.json
	mkdir gen -p &\
	swagger generate model -f swagger.json -t gen

.PHONY: clean
clean: 
	rm -rf gen