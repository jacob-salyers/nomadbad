#! /bin/bash

# Prelude
if [ $CODE_DIR = '' ]
then 
	echo 'please export CODE_DIR' >&2
	exit 1
fi
d=`pwd`
cd $CODE_DIR/nomad

pid=$(sudo lsof -i -P -n | awk '/nomad.*LISTEN/{ print $2; }' | head -1)

sudo kill -9 $pid

# Postlude
cd $d
