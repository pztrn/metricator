# This is a Metricator Makefile.
# It contains calls to scripts placed in scripts directory.

CONFIG ?= ./metricator.example.yaml

help: Makefile
	@echo -e "Metricator Makefile available subcommands:\n"
	@cat $< | grep "## " | sort | sed -n 's/^## //p'
	@echo ""
	@make show-git-data

.DEFAULT_GOAL := help

check-build-dir:
	@if [ ! -d "._bin" ]; then mkdir ._bin; fi

## metricator-client-build: builds metricator client and places into ${PWD}/._bin.
metricator-client-build: check-build-dir
	@if [ -f ./._bin/metricator-client ]; then rm ./._bin/metricator-client; fi
	@scripts/build.sh metricator-client

## metricatord-build: builds metricator daemon and places into ${PWD}/._bin.
metricatord-build: check-build-dir
	@if [ -f ./._bin/metricatord ]; then rm ./._bin/metricatord; fi
	@scripts/build.sh metricatord

## metricator-client-run: starts metricator client. Use ARGS to supply args.
metricator-client-run: metricator-client-build
	@./._bin/metricator-client -config ${CONFIG} $(ARGS)

## metricatord-run: starts metricator daemon.
metricatord-run: metricatord-build
	./._bin/metricatord -config ${CONFIG}

show-git-data:
	@echo "Parameters for current source code state:"
	@scripts/show_git_data.sh
