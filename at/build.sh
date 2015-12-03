#!/bin/sh
TARGET_FILES="at_zmq_proxy at_debug at_jeus at_report at_jmc ck_message if_rs232 if_webconn_rs232 if_webconn_relay"
for FILE in ${TARGET_FILES}
do
  go install -tags zmq_4_x ${FILE}
done
