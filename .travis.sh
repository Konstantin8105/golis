#!/usr/bin/env bash

cd $GOPATH/src/github.com/Konstantin8105/golis/

mkdir bin

LIS_INSTALL_FOLDER=`pwd`
echo "Lis install folder: $LIS_INSTALL_FOLDER"

git clone https://github.com/anishida/lis.git

cd lis

./configure --prefix="$LIS_INSTALL_FOLDER/bin" --enable-quad --enable-omp
# --enable-mpi

make install

cd ..
