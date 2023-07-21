SERVER := chaosmeta
BUILD_DIST := build/dist/
CHAOSMETA_EXEC := ${BUILD_DIST}${SERVER}
NOW := $(shell date "+%Y/%m/%d/%H:%m:%S")

all: ${SERVER}

.PHONY: ${SERVER}
${SERVER}:
	 CGO_ENABLED=0 GOARCH=amd64 go build  -ldflags "-X main.BuildTime=${NOW} -s" -o ${CHAOSMETA_EXEC} cmd/server/main.go


start: ${SERVER}
	./${CHAOSMETA_EXEC} server start --config=conf/app.yaml

air:
	air server start --config=config/app.yaml

.PHONY: ${OBJ}
server: all
	./${SERVER} server

clean:
	-rm ${SERVER}

