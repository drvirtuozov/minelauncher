LAUNCHER=minelauncher
LDFLAGS=-ldflags "-w -s -X main.launcher=${LAUNCHER}"

build:
		go build ${LDFLAGS} -o bin/${LAUNCHER}