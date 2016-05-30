#!/usr/bin/env bash

set -u
set -e

PATH_INSTALL="/opt/coscale/cli"
PATH_SYMLINK="/usr/bin/coscale-cli"
SYMLINK="false"

# Check for unexisting variables
if [ -z ${COSCALE_APPID+x} ]; then
    COSCALE_APPID=""
fi

if [ -z ${COSCALE_TOKEN+x} ]; then
    COSCALE_TOKEN=""
fi

ask_yesno() {
    local message=$1
    while true; do
        read -r -p "$message (yes/no)" yn
        case $yn in
            [Yy]* )
                break
                ;;
            [Nn]* )
                break
                ;;
            * )
                echo "Please answer yes or no."
                ;;
        esac
    done

    if echo "$yn" | grep -iq "^y" ;then
        return 0
    else
        return 1
    fi
}

echo "----------------------------------------------"
echo -e "\tCoScale CLI installation script"
echo "----------------------------------------------"

read -i "$COSCALE_APPID" -p "Please enter your application id: " -e -r COSCALE_APPID

read -i "$COSCALE_TOKEN" -p "Please enter your access token: " -e -r COSCALE_TOKEN

read -i "$PATH_INSTALL" -p "Installation path: " -e -r PATH_INSTALL
PATH_INSTALL=$(readlink -f $PATH_INSTALL || echo $PATH_INSTALL)

if ask_yesno "Create a symlink [$PATH_SYMLINK]?"; then
    SYMLINK="true"
fi
echo

echo "----------------------------------------------"
echo -e "\t Summary"
echo
echo -e "Application ID \t\t: $COSCALE_APPID"
echo -e "Access token \t\t: $COSCALE_TOKEN"
echo
echo -e "Installation path \t: $PATH_INSTALL"
echo -e "Symlink \t\t: $SYMLINK"
if [ "$SYMLINK" = true ] ; then
    echo -e "Symlink path \t\t: $PATH_SYMLINK"
fi
echo
if ! ask_yesno "Is everything correct?"; then
    exit 1
fi

echo "----------------------------------------------"
echo -e "Starting installation"

# Fetch latest release list from Github
echo
echo -e "- Getting latest release information"
github_data=$(curl -s -L https://api.github.com/repos/CoScale/coscale-cli/releases/latest | grep "browser_download_url" | awk '{ print $2; }' | sed 's/"//g')
release=$(echo "$github_data" | grep -v ".exe")

# Create dirs
if [ ! -d "$PATH_INSTALL" ]; then
    echo echo -e "- Creating directories $PATH_INSTALL/"
    mkdir -v -p "$PATH_INSTALL"
    echo
fi

# Install client
echo -e "- Downloading client to $PATH_INSTALL/coscale-cli"
curl -L "$release" > "$PATH_INSTALL/coscale-cli"
chmod -v +x "$PATH_INSTALL/coscale-cli"
echo

# Create symlink from /usr/bin/coscale-cli to /opt/coscale/cli/coscale-cli
if [ "$SYMLINK" = true ] ; then
    echo -e "- Creating symlink"
    # Check if file exists
    if [ -f "/usr/bin/coscale-cli" ]; then
        # Check if symlink is correct
        if [ "$(readlink /usr/bin/coscale-cli)" = "$PATH_INSTALL/coscale-cli" ]; then
            echo -e "\tExisting symlink detected"
        else
            echo -e "\tIncorrect symlink detected, please remove the file /usr/bin/coscale-cli and start again"
            exit 1
        fi
    else
        # Symlink does not exist, create
        ln -v -s "$PATH_INSTALL/coscale-cli" /usr/bin/coscale-cli
        echo -e "\tSymlink created"
    fi
    echo
fi

# Create config
if [ ! -f "$PATH_INSTALL/api.conf" ]; then
    echo -e "- Generating config file"
    echo "{\"baseurl\":\"https://api.coscale.com\", \"appid\":\"$COSCALE_APPID\", \"accesstoken\":\"$COSCALE_TOKEN\"}" | gzip -c > $PATH_INSTALL/api.conf
    echo
fi

# Test config
echo -e "- Testing configuration"
$PATH_INSTALL/coscale-cli check-config | sed -e 's/[{}]//g' | awk -F ":" '{print $2 }'
echo

# Done
echo -e "Done, you can now start using the CoScale CLI tool."
