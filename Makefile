# Variable for filename for store running procees id
PID_FILE = /tmp/packback-be.pid
# We can use such syntax to get main.go and other root Go files.
GO_FILES = $(wildcard *.go)
APP_SERVER = cmd/server/server.go

# Start task performs "go run main.go" command and writes it's process id to PID_FILE.
start:
	go run $(APP_SERVER) & echo $$! > $(PID_FILE)

# Stop task will kill process by ID stored in PID_FILE (and all child processes by pstree).
stop:
	# -kill `pstree -p \`cat $(PID_FILE)\` | tr "\n" " " |sed "s/[^0-9]/ /g" |sed "s/\s\s*/ /g"`
	pkill -9 -P `cat $(PID_FILE)`
# Before task will only prints message. Actually, it is not necessary. You can remove it, if you want.
before:
	@echo "STOPED server" && printf '%*s\n' "40" '' | tr ' ' -

# Restart task will execute stop, before and start tasks in strict order and prints message.
restart: stop before start
	@echo "STARTED server" && printf '%*s\n' "40" '' | tr ' ' -

# Serve task will run fswatch monitor and performs restart task if any source file changed. Before serving it will execute start task.
serve: start
	fswatch -or --event=Updated . | \
	xargs -n1 -I {} make restart

# .PHONY is used for reserving tasks words
.PHONY: start before stop restart serve

kill:
	@echo "Killing process running on port 8000"
	kill $(lsof -t -i:8000)