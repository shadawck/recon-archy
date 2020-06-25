MODULE=recon-archy

docs:
	cd docs && make html

clean:
	rm -rf .vscode

build:
	go build

install:
	go install .

run:
	recon-archy


.PHONY:  clean fclean docs
