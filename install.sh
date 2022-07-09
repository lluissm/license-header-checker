#!/bin/sh

GITHUB_RELEASE=https://github.com/lluissm/license-header-checker/releases/download/v1.3.0
CMD=license-header-checker
FILENAME="${CMD}_linux_amd64"
INSTALL_PATH=/usr/local/bin

if [[ $OSTYPE == 'darwin'* ]]; then
    if [[ $(sysctl -n machdep.cpu.brand_string) =~ "Apple" ]]; then
        FILENAME="${CMD}_macos_arm64"
    else
		FILENAME="${CMD}_macos_intel"
	fi
fi

DOWNLOAD_URL="${GITHUB_RELEASE}/${FILENAME}"
echo "Downloading ${DOWNLOAD_URL}"
curl -LJOs $DOWNLOAD_URL

echo "Installing ${CMD} in ${INSTALL_PATH}"
chmod +x $FILENAME
sudo mv $FILENAME "${INSTALL_PATH}/${CMD}"

$CMD -version
