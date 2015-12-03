#!/bin/sh
TARGET_FILES="server_rs232"
for FILE in ${TARGET_FILES}
do
  go install -tags zmq_4_x ${FILE}
done
