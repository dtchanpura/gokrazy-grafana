VERSION := 9.1.7

all: _gokrazy/extrafiles_arm64.tar _gokrazy/extrafiles_amd64.tar

_gokrazy/extrafiles_amd64.tar:
	mkdir -p _gokrazy/extrafiles_amd64/usr/local/bin
	cp external/grafana/LICENSE _gokrazy/extrafiles_amd64/usr/local/bin/LICENSE.grafana
	cp external/grafana/bin/linux-amd64/grafana-* _gokrazy/extrafiles_amd64/usr/local/bin/
	cd _gokrazy/extrafiles_amd64 && tar cf ../extrafiles_amd64.tar *
	rm -rf _gokrazy/extrafiles_amd64

_gokrazy/extrafiles_arm64.tar:
	mkdir -p _gokrazy/extrafiles_arm64/usr/local/bin
	cp external/grafana/LICENSE _gokrazy/extrafiles_arm64/usr/local/bin/LICENSE.grafana
	cp external/grafana/bin/linux-arm64/grafana-* _gokrazy/extrafiles_arm64/usr/local/bin/
	cd _gokrazy/extrafiles_arm64 && tar cf ../extrafiles_arm64.tar *
	rm -rf _gokrazy/extrafiles_arm64

clean:
	rm -f _gokrazy/extrafiles_*.tar
