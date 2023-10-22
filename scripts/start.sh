#! /bin/bash

# Prelude
if [ $CODE_DIR = '' ]
then 
	echo 'please export CODE_DIR' >&2
	exit 1
fi
d=`pwd`
cd $CODE_DIR/nomad

sudo nohup ./nomad -s 2>&1 &

# Postlude
cd $d
