.PHONY: run.http
run.http:
	go build -o ./build/task-tracker-service ./cmd
	./build/task-tracker-service

.PHONY: run.local
run.local:
	go build -o ./build/task-tracker-service ./cmd
	./build/task-tracker-service -input-method cmd -storage-type local

