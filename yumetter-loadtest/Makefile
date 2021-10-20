.PHONY: run
run:
	YUMETTER_API=https://api.intern-omori.sandbox.yumemi.io sbt "gatling:testOnly sim.BasicSimulation"
	
.PHONY: prof
prof:
	YUMETTER_API=https://api.intern-omori.sandbox.yumemi.io sbt "gatling:testOnly sim.BasicSimulation" & 
	(sleep 30s; go tool pprof --http=:6061 http://10.0.33.47:6060/debug/fgprof?seconds=100)
