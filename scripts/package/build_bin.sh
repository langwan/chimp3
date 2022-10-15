#!/bin/bash

echo "清理"
rm -rf ../../bin
mkdir ../../bin

cd ../../backend
echo "编译mac"
go build -o ../bin/chimp3_backend

echo "编译win"
GOOS=windows GOARCH=amd64 go build -o ../bin/chimp3_backend.exe

echo "编译前端"
cd ../frontend
yarn build