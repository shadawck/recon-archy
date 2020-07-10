MODULE=recon-archy

docs:
	cd docs && make html

# update selenium dependencies
selenium-dep:
	cd ~/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor && go run init.go --alsologtostderr  --download_browsers --download_latest

export:
	/bin/sh -c "export PATH=$$PATH:$$(dirname $$(go list -f '{{.Target}}' .))"

build:
	go build

install:
	go install .

run: build install export
	./recon-archy

# For selenium standalone 
launch: 
	java -jar ~/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/selenium-server-standalone.jar &
	./recon-archy

clean:
	rm -rf .vscode
	rm -rf *.log

fclean: clean
	rm -rf $(MODULE)

.PHONY:  clean fclean docs
