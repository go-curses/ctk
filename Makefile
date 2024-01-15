#!/usr/bin/make -f

SHELL := /bin/bash

BUILD_CMD := go-ctk
DEV_EXAMPLE := demo-app
CDK_PATH := ../cdk
CTK_PATH := ../ctk

.PHONY: all build build-all clean clean-logs dev examples fmt help profile.cpu profile.mem run tidy

all: help

help:
	@echo "usage: make [target]"
	@echo
	@echo "qa targets:"
	@echo "  vet         - run go vet command"
	@echo "  test        - perform all available tests"
	@echo "  cover       - perform all available tests with coverage report"
	@echo
	@echo "cleanup targets:"
	@echo "  clean       - cleans package and built files"
	@echo "  clean-logs  - cleans *.log from the project"
	@echo
	@echo "go.mod helpers:"
	@echo "  local       - add go.mod local CDK package replacements"
	@echo "  unlocal     - remove go.mod local CDK package replacements"
	@echo
	@echo "build targets:"
	@echo "  deps        - install stringer and bitmasker tools"
	@echo "  generate    - run go generate"
	@echo "  examples    - builds all examples"
	@echo "  build       - build the go-ctk command"
	@echo "  build-all   - build all commands"
	@echo "  dev         - build ${DEV_EXAMPLE} with profiling"
	@echo "  *.so        - build a plugin-world shared object file"
	@echo "  *           - build only the given example (by name)"
	@echo
	@echo "run targets:"
	@echo "  run         - run the dev build (sanely handle crashes)"
	@echo "  profile.cpu - run the dev build and profile CPU"
	@echo "  profile.mem - run the dev build and profile memory"

vet:
	@echo -n "# vetting ctk ..."
	@go vet && echo " done"

test: vet
	@echo "# testing ctk ..."
	@go test -v ./...

cover:
	@echo "# testing ctk (with coverage) ..."
	@go test -cover -coverprofile=coverage.out ./...
	@echo "# test coverage ..."
	@go tool cover -html=coverage.out

clean-build-logs:
	@echo "# cleaning *.build.log files"
	@rm -fv *.build.log || true

clean-logs:
	@echo "# cleaning *.log files"
	@rm -fv *.log || true
	@echo "# cleaning *.out files"
	@rm -fv *.out || true
	@echo "# cleaning pprof files"
	@rm -rfv /tmp/*.cdk.pprof || true

clean-cmd: clean-build-logs
	@echo "# cleaning built commands"
	@for tgt in `ls cmd`; do \
		if [ -f $$tgt ]; then rm -fv $$tgt; fi; \
	done

clean-examples: clean-build-logs
	@echo "# cleaning built examples"
	@rm -fv *.so         || true
	@rm -fv hello-plugin || true
	@for tgt in `ls examples`; do \
		if [ -f $$tgt ]; then rm -fv $$tgt; fi; \
	done

clean: clean-logs clean-examples clean-cmd
	@echo "# cleaning goland builds"
	@rm -rfv go_* || true

build:
	@echo -n "# building command ${BUILD_CMD}... "
	@cd cmd/${BUILD_CMD}; \
		( go build -v \
				-trimpath \
				-o ../../${BUILD_CMD} \
			2>&1 ) > ../../${BUILD_CMD}.build.log; \
		rv="$$?"; \
		cd - > /dev/null; \
		if [ $$rv = "0" -a -f ${BUILD_CMD} ]; then \
			echo "done."; \
		else \
			echo "failed.\n>\tsee ./${BUILD_CMD}.build.log for errors"; \
			false; \
		fi

build-all: clean-cmd
	@for tgt in `ls cmd`; \
	do \
		if [ -d cmd/$$tgt ]; \
		then \
			echo -n "# building command $$tgt... "; \
			cd cmd/$$tgt; \
			( go build -v \
					-trimpath \
					-gcflags=all="-N -l" \
					-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
					-o ../../$$tgt \
				2>&1 ) > ../../$$tgt.build.log; \
			cd - > /dev/null; \
			if [ -f $$tgt ]; \
			then \
				echo "done."; \
			else \
				echo "fail.\n#\tsee ./$$tgt.build.log"; \
			fi; \
		fi; \
	done

deps:
	@echo "# installing dependencies..."
	@echo "#\tinstalling stringer"
	@GO111MODULE=off go install golang.org/x/tools/cmd/stringer
	@echo "#\tinstalling bitmasker"
	@GO111MODULE=off go install github.com/go-curses/bitmasker

generate:
	@echo "# generate go sources..."
	@go generate -v ./...

examples: clean-examples hello-plugin.so hello-plugin
	@echo "# building all examples..."
	@for tgt in `ls examples`; \
	do \
		if [ -d examples/$$tgt ]; \
		then \
			echo -n "#\tbuilding example $$tgt... "; \
			cd examples/$$tgt; \
			( go build -v \
					-trimpath \
					-gcflags=all="-N -l" \
					-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
					-o ../../$$tgt \
				2>&1 ) > ../../$$tgt.build.log; \
			cd - > /dev/null; \
			if [ -f $$tgt ]; \
			then \
				echo "done."; \
			else \
				echo "fail.\n#\tsee ./$$tgt.build.log"; \
			fi; \
		fi; \
	done

depends-on-cdk-path:
	@if [ ! -d ${CDK_PATH} ]; then \
			echo "Error: $(MAKECMDGOALS) depends upon a valid CDK_PATH."; \
			echo "Default: ../cdk"; \
			echo ""; \
			echo "Specify the path to an existing CDK checkout with the"; \
			echo "CDK_PATH variable as follows:"; \
			echo ""; \
			echo " make CDK_PATH=../path/to/cdk $(MAKECMDGOALS)"; \
			echo ""; \
			false; \
		fi

tidy:
	@echo "# running go mod tidy"
	@go mod tidy

local: depends-on-cdk-path
	@echo "# adding go.mod local CDK package replacements..."
	@go mod edit -replace=github.com/go-curses/cdk=${CDK_PATH}
	@for tgt in charset encoding env log memphis; do \
		if [ -f ${CDK_PATH}/$$tgt/go.mod ]; then \
			echo "#\t$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/$$tgt=${CDK_PATH}/$$tgt ; \
		fi; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -f ${CDK_PATH}/lib/$$tgt/go.mod ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/lib/$$tgt=${CDK_PATH}/lib/$$tgt ; \
		fi; \
	done

unlocal: depends-on-cdk-path
	@echo "# removing go.mod local CDK package replacements..."
	@go mod edit -dropreplace=github.com/go-curses/cdk
	@for tgt in charset encoding env log memphis; do \
		if [ -f ${CDK_PATH}/$$tgt/go.mod ]; then \
			echo "#\t$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/$$tgt ; \
		fi; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -f ${CDK_PATH}/lib/$$tgt/go.mod ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/lib/$$tgt ; \
		fi; \
	done

dev:
	@if [ -d examples/${DEV_EXAMPLE} ]; \
	then \
		echo -n "# building: ${DEV_EXAMPLE} [dev]... "; \
		cd examples/${DEV_EXAMPLE}; \
		( go build -v \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../${DEV_EXAMPLE} \
			2>&1 ) > ../../${DEV_EXAMPLE}.build.log; \
		rv="$$?"; \
		cd - > /dev/null; \
		if [ $$rv = "0" -a -f ${DEV_EXAMPLE} ]; then \
			echo "done."; \
		else \
			echo "failed.\n>\tsee ./${DEV_EXAMPLE}.build.log errors below:"; \
			cat ./${DEV_EXAMPLE}.build.log; \
			false; \
		fi; \
	else \
		echo "# dev example not found: ${DEV_EXAMPLE}"; \
	fi

run: export GO_CDK_LOG_FILE=./${DEV_EXAMPLE}.cdk.log
run: export GO_CDK_LOG_LEVEL=debug
run: export GO_CDK_LOG_FULL_PATHS=true
run:
	@if [ -f ${DEV_EXAMPLE} ]; \
	then \
		echo "# running: ${DEV_EXAMPLE}"; \
		( ./${DEV_EXAMPLE} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${DEV_EXAMPLE} crashed, see: ./${DEV_EXAMPLE}.cdk.log"; \
			read -p "# reset terminal? [Yn] " RESP; \
			if [ "$$RESP" = "" -o "$$RESP" = "Y" -o "$$RESP" = "y" ]; \
			then \
				reset; \
				echo "# ${DEV_EXAMPLE} crashed, terminal reset, see: ./${DEV_EXAMPLE}.cdk.log"; \
			fi; \
		else \
			echo "# ${DEV_EXAMPLE} exited normally."; \
		fi; \
	fi

profile.cpu: export GO_CDK_LOG_FILE=./${DEV_EXAMPLE}.cdk.log
profile.cpu: export GO_CDK_LOG_LEVEL=debug
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE_PATH=/tmp/${DEV_EXAMPLE}.cdk.pprof
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: dev
	@mkdir -v /tmp/${DEV_EXAMPLE}.cdk.pprof 2>/dev/null || true
	@if [ -f ${DEV_EXAMPLE} ]; \
		then \
			./${DEV_EXAMPLE} && \
			if [ -f /tmp/${DEV_EXAMPLE}.cdk.pprof/cpu.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${DEV_EXAMPLE}.cdk.pprof/cpu.pprof ; \
			else \
				echo "# missing /tmp/${DEV_EXAMPLE}.cdk.pprof/cpu.pprof"; \
			fi ; \
		fi

profile.mem: export GO_CDK_LOG_FILE=./${DEV_EXAMPLE}.log
profile.mem: export GO_CDK_LOG_LEVEL=debug
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE_PATH=/tmp/${DEV_EXAMPLE}.cdk.pprof
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: dev
	@mkdir -v /tmp/${DEV_EXAMPLE}.cdk.pprof 2>/dev/null || true
	@if [ -f ${DEV_EXAMPLE} ]; \
		then \
			./${DEV_EXAMPLE} && \
			if [ -f /tmp/${DEV_EXAMPLE}.cdk.pprof/mem.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${DEV_EXAMPLE}.cdk.pprof/mem.pprof; \
			else \
				echo "# missing /tmp/${DEV_EXAMPLE}.cdk.pprof/mem.pprof"; \
			fi ; \
		fi

%.so: PLUGNAME=$(basename $@)
%.so:
	@if [ -d examples/plugin-world/$(PLUGNAME) ]; \
	then \
		echo -n "# building plugin $(PLUGNAME)... "; \
		cd examples/plugin-world/$(PLUGNAME); \
		( go build -v \
				-buildmode=plugin \
				-trimpath \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../../$@ \
			2>&1 ) > ../../../$(PLUGNAME).build.log; \
		cd - > /dev/null; \
		if [ -f $@ ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./$(PLUGNAME).build.log"; \
		fi; \
	else \
		echo "not a plugin: $@"; \
		false; \
	fi

%:
	@if [ -d examples/$@ ]; \
	then \
		echo -n "# building example $@... "; \
		cd examples/$@; \
		( go build -v \
				-trimpath \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../$@ \
			2>&1 ) > ../../$@.build.log; \
		cd - > /dev/null; \
		if [ -f $@ ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./$@.build.log"; \
		fi; \
	elif [ -d examples/plugin-world/$@ ]; \
	then \
		echo -n "# building example plugin-world/$@... "; \
		cd examples/plugin-world/$@; \
		( go build -v \
				-trimpath \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../../$@ \
			2>&1 ) > ../../../$@.build.log; \
		cd - > /dev/null; \
		if [ -f $@ ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./$@.build.log"; \
		fi; \
	elif [ -d cmd/$@ ]; \
	then \
		echo -n "# building command $@... "; \
		cd cmd/$@; \
		( go build -v \
				-trimpath \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../$@ \
			2>&1 ) > ../../$@.build.log; \
		cd - > /dev/null; \
		if [ -f $@ ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./$@.build.log"; \
		fi; \
	else \
		echo "not a command or example: $@"; \
	fi

#
# Cross-Compilation Targets
#

dev-linux-mips64: export GOOS=linux
dev-linux-mips64: export GOARCH=mips64
dev-linux-mips64: dev

dev-darwin-amd64: export GOOS=darwin
dev-darwin-amd64: export GOARCH=amd64
dev-darwin-amd64: dev
