# This is a Metricator Makefile.
# It contains calls to scripts placed in scripts directory.

help: Makefile
	@echo -e "Metricator Makefile available subcommands:\n"
	@cat $< | grep "## " | sort | sed -n 's/^## //p'
	@echo ""
	@make -f Makefile show-git-data

.DEFAULT_GOAL := help

check-build-dir:
	@if [ ! -d "._bin" ]; then mkdir ._bin; fi

## metricatord-build: builds metricator daemon and places into ${PWD}/._bin.
metricatord-build: check-build-dir
	rm ./._bin/metricatord || true
	cd cmd/metricatord && go build -o ../../._bin/metricatord

## metricatord-run: starts metricator daemon.
metricatord-run: metricatord-build
	./._bin/metricatord -config ./metricator.example.yaml

show-git-data:
	@echo "Parameters for current source code state:"
	@scripts/show_git_data.sh
