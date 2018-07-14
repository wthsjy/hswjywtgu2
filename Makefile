tag_version=v0.27-beta
tag_desc="make tag test"

tag:
	go fmt ./...
	git add .
	git commit -m $(tag_desc)
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	git push
	git tag $(tag_version)
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	git push origin $(tag_version)

dep_init:
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	dep init -v