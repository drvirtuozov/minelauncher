LAUNCHER=minelauncher
MINEVER=1.12.2
ASSET_INDEX=1.12
CLIENT_URL=https://bitbucket.org/drvirtuozov/mineclient/get/master.zip
CLIENT_TOKEN=78cbfcf2-b52c-4420-bf7c-3e569855d6e5
LDFLAGS=-ldflags "-w -s\
		-X main.launcher=${LAUNCHER}\
		-X main.minever=${MINEVER}\
		-X main.assetIndex=${ASSET_INDEX}\
		-X main.clientURL=${CLIENT_URL}\
		-X main.clientToken=${CLIENT_TOKEN}\
"

build:
		go build ${LDFLAGS} -o bin/${LAUNCHER}