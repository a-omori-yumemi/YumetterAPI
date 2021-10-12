swag.gen: swagger.json
	mkdir -p gen &\
	swagger generate model -f swagger.json -t gen

.PHONY: clean
clean: 
	rm -rf gen