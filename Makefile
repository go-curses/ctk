#!/usr/bin/make -f

BUILD_CMD := go-ctk
DEV_EXAMPLE := demo-app
CDK_PATH := ../cdk

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
	@echo "  examples    - builds all examples"
	@echo "  build       - build the go-ctk command"
	@echo "  build-all   - build all commands"
	@echo "  dev         - build ${DEV_EXAMPLE} with profiling"
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
	@for tgt in `ls examples`; do \
		if [ -f $$tgt ]; then rm -fv $$tgt; fi; \
	done

clean: clean-logs clean-examples clean-cmd
	@echo "# cleaning goland builds"
	@rm -rfv go_* || true

build: clean-cmd
	@echo -n "# building command ${BUILD_CMD}... "
	@cd cmd/${BUILD_CMD}; \
		( go build -v \
				-o ../../${BUILD_CMD} \
			2>&1 ) > ../../${BUILD_CMD}.build.log; \
		cd ../..; \
		if [ -f ${BUILD_CMD} ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./${BUILD_CMD}.build.log"; \
		fi

build-all: clean-cmd
	@for tgt in `ls cmd`; \
	do \
		if [ -d cmd/$$tgt ]; \
		then \
			echo -n "# building command $$tgt... "; \
			cd cmd/$$tgt; \
			( go build -v \
					-tags `echo "cmd-$$tgt" | perl -pe 's/-/_/g'` \
					-o ../../$$tgt \
				2>&1 ) > ../../$$tgt.build.log; \
			cd ../..; \
			if [ -f $$tgt ]; \
			then \
				echo "done."; \
			else \
				echo "fail.\n#\tsee ./$$tgt.build.log"; \
			fi; \
		fi; \
	done

examples: clean-examples
	@echo "# building all examples..."
	@for tgt in `ls examples`; \
	do \
		if [ -d examples/$$tgt ]; \
		then \
			echo -n "#\tbuilding example $$tgt... "; \
			cd examples/$$tgt; \
			( go build -v \
					-tags `echo "debug example-$$tgt" | perl -pe 's/-/_/g'` \
					-o ../../$$tgt \
				2>&1 ) > ../../$$tgt.build.log; \
			cd ../..; \
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

local: depends-on-cdk-path
	@echo "# adding go.mod local CTK package replacements..."
	@go mod edit -replace=github.com/go-curses/ctk=../ctk
	@echo "# adding go.mod local CDK package replacements..."
	@go mod edit -replace=github.com/go-curses/cdk=${CDK_PATH}
	@for tgt in charset encoding env log memphis; do \
			echo "#\t$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/$$tgt=${CDK_PATH}/$$tgt ; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -d ${CDK_PATH}/lib/$$tgt ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/lib/$$tgt=${CDK_PATH}/lib/$$tgt ; \
		fi; \
	done

unlocal: depends-on-cdk-path
	@echo "# removing go.mod local CTK package replacements..."
	@go mod edit -dropreplace=github.com/go-curses/ctk
	@echo "# removing go.mod local CDK package replacements..."
	@go mod edit -dropreplace=github.com/go-curses/cdk
	@for tgt in charset encoding env log memphis; do \
			echo "#\t$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/$$tgt ; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -d ${CDK_PATH}/lib/$$tgt ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/lib/$$tgt ; \
		fi; \
	done

dev: clean
	@if [ -d examples/${DEV_EXAMPLE} ]; \
	then \
		echo -n "# building: ${DEV_EXAMPLE} [dev]... "; \
		cd examples/${DEV_EXAMPLE}; \
		( go build -v \
				-tags `echo "debug example-${DEV_EXAMPLE}" | perl -pe 's/-/_/g'` \
				-ldflags="-X 'main.IncludeProfiling=true'" \
				-gcflags=all="-N -l" \
				-o ../../${DEV_EXAMPLE} \
			2>&1 ) > ../../${DEV_EXAMPLE}.build.log; \
		cd ../..; \
		[ -f ${DEV_EXAMPLE} ] \
			&& echo "done." \
			|| echo "failed.\n>\tsee ./${DEV_EXAMPLE}.build.log for errors"; \
	else \
		echo "# dev example not found: ${DEV_EXAMPLE}"; \
	fi

run: export GO_CDK_LOG_FILE=./${DEV_EXAMPLE}.cdk.log
run: export GO_CDK_LOG_LEVEL=debug
run: export GO_CDK_LOG_FULL_PATHS=true
run: dev
	@if [ -f ${DEV_EXAMPLE} ]; \
	then \
		echo "# starting ${DEV_EXAMPLE}..."; \
		./${DEV_EXAMPLE}; \
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

%:
	@if [ -f $@ -o -f $@.build.log ]; \
	then \
		echo -n "# cleaning $@... "; \
		rm -f $@ $@.build.log; \
		echo "done."; \
	fi; \
	if [ -d examples/$@ ]; \
	then \
		echo -n "# building example $@... "; \
		cd examples/$@; \
		( go build -v \
				-tags `echo "debug example-$@" | perl -pe 's/-/_/g'` \
				-o ../../$@ \
			2>&1 ) > ../../$@.build.log; \
		cd ../..; \
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
				-tags `echo "cmd-$@" | perl -pe 's/-/_/g'` \
				-o ../../$@ \
			2>&1 ) > ../../$@.build.log; \
		cd ../..; \
		if [ -f $@ ]; \
		then \
			echo "done."; \
		else \
			echo "fail.\n#\tsee ./$@.build.log"; \
		fi; \
	else \
		echo "not a command or example: $@"; \
	fi
