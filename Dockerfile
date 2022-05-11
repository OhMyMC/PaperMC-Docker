FROM azul/zulu-openjdk-alpine:17-jre AS build
RUN apk add curl jq

WORKDIR /opt/minecraft
COPY scripts/getpaperserver.sh /
COPY .env /

RUN chmod +x /getpaperserver.sh
RUN /getpaperserver.sh

FROM azul/zulu-openjdk-alpine:17-jre AS runtime

LABEL maintainer="MoFan <suixinio@163.com>"

ENV MINECRAFT_PATH=/opt/minecraft
ENV MINECRAFT_DATA_PATH=/data
ENV MEMORY_SIZE=2G
ENV JAVA_ARGS="-Dcom.mojang.eula.agree=true"
ENV PAPERMC_ARG="--nojline nogui"

WORKDIR ${MINECRAFT_DATA_PATH}

COPY --from=build ${MINECRAFT_PATH}/papermc.jar ${MINECRAFT_PATH}/papermc.jar

WORKDIR ${MINECRAFT_DATA_PATH}

COPY /scripts/docker-entrypoint.sh /opt/minecraft
RUN chmod +x /opt/minecraft/docker-entrypoint.sh

RUN addgroup minecraft
RUN adduser -s /bin/bash minecraft -G minecraft -h ${MINECRAFT_PATH} -D

RUN chown -R minecraft:minecraft ${MINECRAFT_PATH}
RUN chown -R minecraft:minecraft ${MINECRAFT_DATA_PATH}

USER minecraft

VOLUME "${MINECRAFT_DATA_PATH}"

# Expose minecraft port
EXPOSE 25565/tcp
EXPOSE 25565/udp

# Entrypoint
ENTRYPOINT ["/opt/minecraft/docker-entrypoint.sh"]
