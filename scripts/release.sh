#!/usr/bin/env bash

usage() {
    echo "USAGE: ./release.sh [version] [msg...]"
    exit 1
}

REVISION=$(git rev-parse HEAD)
GIT_TAG=$(git name-rev --tags --name-only $REVISION)
if [ "$GIT_TAG" = "" ]; then
    GIT_TAG="devel"
fi


VERSION=$1
if [ "$VERSION" = "" ]; then
    echo "Need to specify a version! Perhaps '$GIT_TAG'?"
    usage
fi

set -u -e

rm -rf /tmp/dorepl_build/

mkdir -p /tmp/dorepl_build/linux
GOOS=linux go build -ldflags "-X main.version=$VERSION" -o /tmp/dorepl_build/linux/dorepl github.com/aybabtme/godotto/cmd/dorepl
pushd /tmp/dorepl_build/linux/
tar cvzf /tmp/dorepl_build/dorepl_linux.tar.gz dorepl
popd

mkdir -p /tmp/dorepl_build/darwin
GOOS=darwin go build -ldflags "-X main.version=$VERSION" -o /tmp/dorepl_build/darwin/dorepl github.com/aybabtme/godotto/cmd/dorepl
pushd /tmp/dorepl_build/darwin/
tar cvzf /tmp/dorepl_build/dorepl_darwin.tar.gz dorepl
popd

mkdir -p /tmp/dorepl_build/windows
GOOS=windows go build -ldflags "-X main.version=$VERSION" -o /tmp/dorepl_build/windows/dorepl.exe github.com/aybabtme/godotto/cmd/dorepl
pushd /tmp/dorepl_build/windows/
zip /tmp/dorepl_build/dorepl_windows.zip dorepl.exe
popd


temple file < README.tmpl.md > ../README.md -var "version=$VERSION"
git add ../README.md
git commit -m 'release bump'

hub release create \
    -a /tmp/dorepl_build/dorepl_linux.tar.gz \
    -a /tmp/dorepl_build/dorepl_darwin.tar.gz \
    -a /tmp/dorepl_build/dorepl_windows.zip \
    $VERSION

git push origin master
