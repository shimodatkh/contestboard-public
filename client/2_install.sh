#!/bin/bash

sudo apt update
sudo apt install -y jq
sudo apt install -y unzip


# 直接リポジトリにおいておくのでダウンロードは省略
# wget https://github.com/tkuchiki/alp/releases/download/v1.0.3/alp_linux_amd64.zip
# unzip alp_linux_amd64.zip
# rm alp_linux_amd64.zip
# mv alp ../bin
chmod +x ../bin/alp

# ISUCON本番での環境のGoに干渉しないようにOFF
# wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
# sudo tar -C /usr/local -xvzf go1.16.5.linux-amd64.tar.gz
# export  PATH=$PATH:/usr/local/go/bin
go version

#ubuntu18
# wget https://downloads.percona.com/downloads/percona-toolkit/3.3.1/binary/debian/bionic/x86_64/percona-toolkit_3.3.1-1.bionic_amd64.deb

# #ubuntu20
# wget https://downloads.percona.com/downloads/percona-toolkit/3.3.1/binary/debian/focal/x86_64/percona-toolkit_3.3.1-1.focal_amd64.deb

#sudo apt-get install libdbd-mysql-perl libdbi-perl libio-socket-ssl-perl libnet-ssleay-perl libterm-readkey-perl -y
#sudo dpkg -i percona-toolkit_3.3.1-1.bionic_amd64.deb
#sudo dpkg -i percona-toolkit_3.3.1-1.focal_amd64.deb 

# server only
#sudo apt install -y gcc
