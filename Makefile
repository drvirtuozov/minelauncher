LAUNCHER=minelauncher
MINECRAFT_VERSION=1.12.2
ASSET_INDEX=1.12
CLIENT_URL=https://bitbucket.org/drvirtuozov/mineclient/get/master.zip
LDFLAGS=-ldflags "-w -s -X main.lname=${LAUNCHER} -X main.mversion=${MINECRAFT_VERSION} \
		-X main.assetIndex=${ASSET_INDEX} -X main.clientURL=${CLIENT_URL}"

build:
		go build ${LDFLAGS} -o bin/${LAUNCHER}