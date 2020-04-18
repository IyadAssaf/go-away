#/bin/sh

# Modified from https://raw.githubusercontent.com/getlantern/systray/master/example/icon/make_icon.sh

go get github.com/cratonica/2goarray
if [ $? -ne 0 ]; then
    echo failed to install 2goarray
    exit
fi

: ${GOPATH:=$HOME/go}

echo "Got GOPATH $GOPATH"

ICON=$1
OUTPUT=$2
VARIABLE_NAME=$3

echo generating $OUTPUT
echo "//+build linux darwin" > $OUTPUT
echo >> $OUTPUT
cat $ICON | $GOPATH/bin/2goarray $VARIABLE_NAME main >> $OUTPUT
if [ $? -ne 0 ]; then
    echo failured to generate $OUTPUT
    exit
fi
echo generated $OUTPUT
