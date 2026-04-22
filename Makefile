server:
	watchmedo auto-restart --patterns="*.go" --recursive -- go run cmd/surl/main.go -c config.yaml