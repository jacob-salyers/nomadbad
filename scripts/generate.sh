#! /bin/bash

# Prelude
d=`pwd`
cd ~/code/nomad

for file in `ls templates/pages`
do
sed "/~BODY~/ {
	r templates/pages/$file
	d
}" templates/wrapper.html > static/$file
done

# Postlude
cd $d
