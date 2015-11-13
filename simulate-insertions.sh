#!/bin/bash
while [ 1 == 1 ]; do
  VALUE=`date`
  echo "db.times.insert({\"time\": \"$VALUE\"})" | mongo --port 26016 2> /dev/null > /dev/null
  echo "db.times.find({\"time\": /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/})" | mongo --port 26016 2> /dev/null > /dev/null
  sleep 0.1
done
