#/bin/sh

# Taken from https://raw.githubusercontent.com/getlantern/systray/master/example/icon/make_icon.sh

go get github.com/cratonica/2goarray
if [ $? -ne 0 ]; then
    echo Failure executing go get github.com/cratonica/2goarray
    exit
fi

if [ -z "$1" ]; then
    echo Please specify a PNG file
    exit
fi

if [ ! -f "$1" ]; then
    echo $1 is not a valid file
    exit
fi

OUTPUT=./cmd/taskbar/iconunix.go
echo Generating $OUTPUT
echo "//+build linux darwin" > $OUTPUT
echo >> $OUTPUT
cat "$1" | $GOPATH/bin/2goarray iconData main >> $OUTPUT
if [ $? -ne 0 ]; then
    echo Failure generating $OUTPUT
    exit
fi
echo Finished