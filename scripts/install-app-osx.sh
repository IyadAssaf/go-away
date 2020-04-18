#!/bin/bash

set -e

# OSX package structure
#SystrayApp.app/
#  Contents/
#    Info.plist
#    MacOS/
#      go-executable
#    Resources/
#      SystrayApp.icns

rm -rf ./GoAway.app || true

mkdir -p GoAway.app/Contents/MacOS/
mkdir -p GoAway.app/Contents/Resources/
touch GoAway.app/Contents/Resources/GoAway.icns

cp ./assets/Info.plist.xml GoAway.app/Contents/Info.plist

./scripts/generate-icon.sh ./assets/camera-off.png ./cmd/taskbar/icon_camera_off.go cameraOffIconData
./scripts/generate-icon.sh ./assets/camera-on.png ./cmd/taskbar/icon_camera_on.go cameraOnIconData

GO111MODULE=on
go build -o goaway ./cmd/taskbar
chmod +x goaway
mv goaway GoAway.app/Contents/MacOS/

echo "Built $PWD/GoAway.app, move it to your Applications folder"