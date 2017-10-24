LAUNCHER=minelauncher
MINEVER=1.12.2
ASSET_INDEX=1.12
CLIENT_URL=https://bitbucket.org/drvirtuozov/mineclient/get/master.zip
LDFLAGS=-ldflags "-w -s -X main.launcher=${LAUNCHER} -X main.minever=${MINEVER} -X main.assetIndex=${ASSET_INDEX} -X main.clientURL=${CLIENT_URL}"

build:
		go build ${LDFLAGS} -o bin/${LAUNCHER}