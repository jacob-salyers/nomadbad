#! /bin/bash

# Prelude
if [ $CODE_DIR = '' ]
then 
	echo 'please export CODE_DIR' >&2
	exit 1
fi
d=`pwd`
cd $CODE_DIR/nomad

curl https://calendar.google.com/calendar/ical/d4d0773546ae6d7e69b104a3f0a6e612ee49bf31fcd87463de373f9819be336c%40group.calendar.google.com/private-1250974c781dffb16a03c7b44773371d/basic.ics \
	| node scripts/calendarToHTML.mjs > templates/pages/schedule.html

for file in `ls templates/pages`
do
sed "
1a \
	<!-- DO NOT EDIT THIS FILE -->
/~BODY~/ {
	r templates/pages/$file
	d
}" templates/wrapper.html > static/$file
done

# Postlude
cd $d
