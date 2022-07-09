#!/bin/sh

CMD=license-header-checker
INSTALL_PATH=/usr/local/bin

if [ -z "$1" ]; then
	LATEST=$(curl -fsSLI -o /dev/null -w "%{url_effective}\n" https://github.com/lluissm/license-header-checker/releases/latest | sed 's#.*tag/##g' && echo)
	echo "Installing latest version: ${LATEST}"
	VERSION=$LATEST
else
	VERSION=$1
fi

GITHUB_RELEASE="https://github.com/lluissm/license-header-checker/releases/download/${VERSION}"

FILENAME="${CMD}_linux_amd64"
if [[ $OSTYPE == 'darwin'* ]]; then
	if [[ $(sysctl -n machdep.cpu.brand_string) =~ "Apple" ]]; then
		FILENAME="${CMD}_macos_arm64"
	else
		FILENAME="${CMD}_macos_intel"
	fi
fi

DOWNLOAD_URL="${GITHUB_RELEASE}/${FILENAME}"
echo "Downloading ${DOWNLOAD_URL}"
if ! curl -LJOs --fail "${DOWNLOAD_URL}" -o /dev/null; then
	echo "Could not download the file. Check if the version ${VERSION} exists."
	exit 1
fi

echo "Installing ${CMD} in ${INSTALL_PATH}"
chmod +x $FILENAME
sudo mv $FILENAME "${INSTALL_PATH}/${CMD}"

$CMD -version
