#!/usr/bin/env bash

echo "Preparing to install CoScale CLI"

# Check command arguments
if [ -z "$COSCALE_APPID" ]; then
    echo "App id:"
    read -e COSCALE_APPID
fi

if [ -z "$COSCALE_TOKEN" ]; then
    echo "Access token:"
    read -e COSCALE_TOKEN
fi

# Fetch latest release list from Github
github_data=`curl -s -L https://api.github.com/repos/CoScale/coscale-cli/releases/latest | python -c 'import json,sys;obj=json.load(sys.stdin); releases = [release["browser_download_url"] for release in obj["assets"]]; print "\n".join(releases)'`

# Detect operation system
os=`uname -o | awk '{split($0,a,"/"); print tolower(a[2])}'`

# Select correct release
release=`echo "$github_data" | grep $os`

# Create dirs
mkdir -p /opt/coscale/cli
pushd /opt/coscale/cli

# Install client
curl -L "$release" > coscale-cli
chmod +x coscale-cli

# Add to $PATH
echo "export PATH=\"${PATH}:/opt/coscale/cli\"" >> /root/.bashrc
source /root/.bashrc

# Create config
echo "{\"baseurl\":\"https://api.coscale.com\", \"appid\":\"$COSCALE_APPID\", \"accesstoken\":\"$COSCALE_TOKEN\"}" | gzip -c > /opt/coscale/cli/api.conf

# Test config
coscale-cli event list
