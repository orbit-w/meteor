pwd:=$(shell pwd)

#make build_tag version=<版本> desc="<描述>"
build_tag:
	git tag -a $(version) -m $(desc)
	git push origin $(version)