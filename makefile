build:
	go build -o bo-ana && ./bo-ana
	sudo ln -s ~/project/go/src/boilerplate-generator/bo-ana /usr/local/bin
