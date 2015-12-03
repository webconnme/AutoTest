#!/bin/sh
TARGET_FILES="ap_rs232 ap_relay"
for FILE in ${TARGET_FILES}
do
  go install ${FILE}
done
