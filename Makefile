GO=go
GOFMT=gofmt
GO_BINDATA=go-bindata
OUT=./strava-ical

$(OUT): */**.go utils/assets.go
	$(GO) get
	$(GOFMT) -w */**.go
	$(GO) build

run: $(OUT)
	sudo $(OUT) -c config.toml

utils/assets.go: views/** static/**
	$(GO_BINDATA) -o utils/assets.go -pkg utils views/... static/...

.PHONY: clean

clean:
	rm $(OUT) utils/assets.go

install: $(OUT)
	[ -d /etc/strava-ical/ ] || mkdir /etc/strava-ical/
	[ -d /var/log/strava-ical/ ] || mkdir /var/log/strava-ical/
	cp $(OUT) /usr/local/bin/strava-ical
	[ -e /etc/strava-ical/config.toml ] || cp ./config.example.toml /etc/strava-ical/config.toml
	[ -e /etc/supervisor/conf.d/strava-ical.conf ] || cp ./supervisor.example.conf /etc/supervisor/conf.d/strava-ical.conf