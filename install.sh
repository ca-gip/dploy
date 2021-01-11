#!/bin/bash

checksum () {
   curl -s https://api.github.com/repos/ca-gip/dploy/releases/latest \
  | grep browser_download_url \
  | grep checksums \
  | cut -d '"' -f 4 \
  | xargs curl -sL
}

if [ "$(uname)" == "Darwin" ]; then
  echo "Downloading Darwin Release"
  mkdir -p /var/tmp/dploy
  curl -s https://api.github.com/repos/ca-gip/dploy/releases/latest \
    | grep browser_download_url \
    | grep darwin_amd64 \
    | cut -d '"' -f 4 \
    | xargs curl -sL \
    | tar xf - -C /var/tmp/dploy/
    sudo sh -c 'mv /var/tmp/dploy/dploy /usr/local/bin/ && chmod +x /usr/local/bin/dploy'
    rm -rf /var/tmp/dploy
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  echo "Downloading Linux Release"
  mkdir -p /tmp/dploy
  curl -s  https://api.github.com/repos/ca-gip/dploy/releases/latest \
    | grep browser_download_url \
    | grep linux_amd64 \
    | cut -d '"' -f 4 \
    | xargs curl -sL \
    | tar xzf - -C /tmp/dploy
    sudo sh -c 'mv /tmp/dploy/dploy /usr/local/bin/ && chmod +x /usr/local/bin/dploy'
    rm -rf /tmp/dploy
else echo "Unsupported OS" && exit 1
fi

echo "Install done !"