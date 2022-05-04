#!/bin/sh
set -e
exec java -server -jar -Xms$MEMORY_SIZE -Xmx$MEMORY_SIZE $JAVA_ARGS /opt/minecraft/papermc.jar $PAPERMC_ARG