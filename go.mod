module github.com/bryanl/k8s-lab

go 1.13

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190910110746-680d30ca3117 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.1
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	google.golang.org/grpc v1.23.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gotest.tools v2.2.0+incompatible // indirect
	sigs.k8s.io/kind v0.5.1
)

replace github.com/docker/docker v1.13.1 => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4
