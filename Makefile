tag_version=v0.21

tag:
	git add .
	git commit -m "init test"
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	git push
	git tag $(tag_version)
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	git push origin $(tag_version)

dep_init:
	export http_proxy=http://127.0.0.1:1087;export https_proxy=http://127.0.0.1:1087;
	dep init -v