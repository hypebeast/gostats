#!/bin/sh

####
## Builds release packages for the webserver, crawler and the sync script.
## These packages are used by Ansible to deploy the applications to a remote server.
##
## Part of the GoStats project.
##
## Copyright <2014>, Sebastian Ruml <sebastian.ruml@gmail.com>
##
## Created: 2014.08.17
####


VERSION=`cat ../version`
BUILD_DIR=../build

SRCBIN_CRAWLER=crawler
SRCBINPATH_CRAWLER=../crawler
TGZDIR_CRAWLER=$BUILD_DIR/tgz/crawler

SRCBIN_SERVER=webserver
SRCBINPATH_SERVER=../web
TGZDIR_SERVER=$BUILD_DIR/tgz/webserver

SRCBIN_BIGQUERY=sync.sh
SRCBINPATH_BIGQUERY=../bigquery
TGZDIR_BIGQUERY=$BUILD_DIR/tgz/bigquery

# Create build directories
echo "## Creating build directories"

rm -rf $BUILD_DIR
mkdir -p $TGZDIR_CRAWLER
mkdir -p $TGZDIR_SERVER
mkdir -p $TGZDIR_SERVER
mkdir -p $TGZDIR_BIGQUERY

# Build crawler packages
echo "## Building crawler..."

pushd $SRCBINPATH_CRAWLER
gox -os="linux"

cp crawler_linux_amd64 $TGZDIR_CRAWLER/$SRCBIN_CRAWLER
rm crawler_*
popd

pushd $TGZDIR_CRAWLER/../
tar czf crawler.tar.gz crawler/
popd

echo "## Building crawler... DONE"

# Build webserver package
echo "## Building webserver..."

pushd $SRCBINPATH_SERVER
gox -os="linux" -output="webserver_{{.OS}}_{{.Arch}}"

cp webserver_linux_amd64 $TGZDIR_SERVER/$SRCBIN_SERVER
rm webserver_*

cp -r public $TGZDIR_SERVER
cp -r views $TGZDIR_SERVER
cp start.sh $TGZDIR_SERVER
popd

pushd $TGZDIR_SERVER
find public/javascript ! -name "app.min.js" -type f -delete
find . -name "*.scss" -type f -delete
popd

pushd $TGZDIR_SERVER/../
tar czf webserver.tar.gz webserver/
popd

echo "## Building webserver...DONE"

echo "## Building sync script package..."

pushd $SRCBINPATH_BIGQUERY
cp sync.sh $TGZDIR_BIGQUERY/$SRCBIN_BIGQUERY
popd

pushd $TGZDIR_BIGQUERY/../
tar czf bigquery.tar.gz bigquery/
popd


echo "## Building sync script package...DONE"

echo
echo "#### Done building packages!"
exit 0
