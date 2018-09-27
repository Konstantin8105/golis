#!/usr/bin/env bash

cd $GOPATH/src/github.com/Konstantin8105/golis/

git clone https://github.com/anishida/lis.git

cd lis

./configure --enable-quad

make install

cd ..
