#!/bin/sh

source /.env

Version=$PAPER_BUILD_VERSION
Build=$PAPER_BUILD_NUMBER
curl -H "Accept-Encoding: identity" -H "Accept-Language: en" -L -A "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4.212 Safari/537.36" -o papermc.jar "https://papermc.io/api/v2/projects/paper/versions/$Version/builds/$Build/downloads/paper-$Version-$Build.jar"

echo "Downloaded Paper $Version Build $Build"
