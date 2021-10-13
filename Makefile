swag.gen: swagger.json
	mkdir -p handler/gen &\
	swagger generate model -f swagger.json -t handler/gen

.PHONY: clean
clean: 
	rm -rf gen
