.PHONY: start-server stop-server

start-server:
	@nohup python3 -m http.server 8080 --directory server8080 > server8080.log 2>&1 & echo $$! > server8080.pid
	@nohup python3 -m http.server 8081 --directory server8081 > server8081.log 2>&1 & echo $$! > server8081.pid
	@nohup python3 -m http.server 8082 --directory server8082 > server8082.log 2>&1 & echo $$! > server8082.pid
	@echo "Servers started on ports 8080, 8081, 8082"

stop-server:
	-@if [ -f server8080.pid ]; then kill `cat server8080.pid` 2>/dev/null || true; rm -f server8080.pid; fi
	-@if [ -f server8081.pid ]; then kill `cat server8081.pid` 2>/dev/null || true; rm -f server8081.pid; fi
	-@if [ -f server8082.pid ]; then kill `cat server8082.pid` 2>/dev/null || true; rm -f server8082.pid; fi
	@echo "Servers stopped"