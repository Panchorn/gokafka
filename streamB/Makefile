build:
	cp -R ../encryption .
	cp -R ../events .
	cp -R ../logs .
	docker build -t gokafka-stream-b .
	rm -R encryption/
	rm -R events/
	rm -R logs/
