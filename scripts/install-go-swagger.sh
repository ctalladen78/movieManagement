#! /bin/bash

download_url="https://github.com/go-swagger/go-swagger/releases/download/v0.23.0/swagger_linux_amd64"
sudo curl -o /usr/local/bin/swagger -L'#' "$download_url"
sudo chmod +x /usr/local/bin/swagger
