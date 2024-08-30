build:
	cp -R ../encryption .
	cp -R ../events .
	cp -R ../logs .
	docker build -t gokafka-api-a --progress=plain . &> build.log
	rm -R encryption/
	rm -R events/
	rm -R logs/
