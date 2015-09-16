#!/usr/bin/env bash


#
# Show url after getting latest version
# Improve output formatting
#

set -u
set -e

echo
echo "## Checking configuration and environment"
echo

if [ -z ${COSCALE_APPID+x} ] || [ -z ${COSCALE_TOKEN+x} ]; then
    echo "### Configuration"
    echo
    # Check command arguments
    if [ -z ${COSCALE_APPID+x} ]; then
        echo "Please enter your app id:"
        read -e COSCALE_APPID
    fi
    echo

    if [ -z ${COSCALE_TOKEN+x} ]; then
        echo "Please enter your access token:"
        read -e COSCALE_TOKEN
    fi
    echo
fi

# Detect operation system
echo "### Detecting operating system"
echo
os=`uname -o | awk '{split($0,a,"/"); print tolower(a[2])}'`
echo "Operation system: {{$os}}"
echo

# Fetch latest release list from Github
echo
echo "## Getting latest release information"
echo
github_data=`curl -s -L https://api.github.com/repos/CoScale/coscale-cli/releases/latest | grep "browser_download_url" | awk '{ print $2; }' | sed 's/"//g'`

# Select correct release
release=`echo "$github_data" | grep $os`
echo "### Latest release: ${{release}}"

# Start install
echo
echo "## Installing CoScale CLI tool"
echo

# Create dirs
echo "### Creating directories /opt/coscale/cli"
echo
mkdir -v -p /opt/coscale/cli
pushd /opt/coscale/cli
echo

# Install client
echo "### Downloading client to /opt/coscale/cli/coscale-cli"
curl -L "$release" > coscale-cli
chmod -v +x coscale-cli
echo

# Create symlink from /usr/bin/coscale-cli to /opt/coscale/cli/coscale-cli
echo "### Creating symlink"
ln -v -S /usr/bin/coscale-cli /opt/coscale/cli/coscale-cli
echo

# Create config
echo "Generating config"
echo "{\"baseurl\":\"https://api.coscale.com\", \"appid\":\"$COSCALE_APPID\", \"accesstoken\":\"$COSCALE_TOKEN\"}" | gzip -c > /opt/coscale/cli/api.conf
echo

# Test config
echo "Testing configuration"
coscale-cli event list
echo

# Done
echo "Done, you can now start using the CoScale CLI tool."
popd