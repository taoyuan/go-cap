PREFIX=/usr/local
DESTDIR=
GOFLAGS=
BINDIR=${PREFIX}/bin

BLDDIR = build
EXT=
ifeq (${GOOS},windows)
    EXT=.exe
endif

APPS = capd
all: $(APPS)

$(BLDDIR)/capd:        $(wildcard apps/capd/*.go       cap/*.go       config/*.go internal/*/*.go)

$(BLDDIR)/%:
	@mkdir -p $(dir $@/$*)
	#go build ${GOFLAGS} -o $@ ./apps/$*
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GOFLAGS} -o $@/$*_linux_amd64 ./apps/$*
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build ${GOFLAGS} -o $@/$*_linux_arm ./apps/$*

$(APPS): %: $(BLDDIR)/%

clean:
	rm -fr $(BLDDIR)

.PHONY: clean all
.PHONY: $(APPS)
