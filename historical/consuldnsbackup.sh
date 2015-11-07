#!/bin/bash

export TOKEN="<%= @token %>"
SERVICES_FILE="/etc/consul-template/output/consul-services-generated.ini"
SERVICES_LAST="$SERVICES_FILE.last"
NEW_SHA=$(sha256sum $SERVICES_FILE)
logger -t consul consuldnsbackup: Created new file: $NEW_SHA

if [ -f $SERVICES_FILE ]; then
  if [ -f $SERVICES_LAST ]; then
    OLD_SHA=$(sha256sum $SERVICES_LAST)
    SERVICES_DIFF=$(diff -u $SERVICES_LAST $SERVICES_FILE)
    curl -s -X POST -H "Content-type: application/json" \
    -d "{
          \"title\": \"Consul Services Update\",
          \"text\": \"New File: $NEW_SHA\nOld File: $OLD_SHA\n\nDiff:$SERVICES_DIFF\",
          \"alert_type\": \"info\"
      }" \
    '<%= @url %>/api/v1/events?api_key=<%= @api_key %>' > /dev/null
  fi
  cp -f $SERVICES_FILE $SERVICES_LAST
  SERVICES_DATA=$(cat $SERVICES_FILE)
  /usr/local/bin/consulkv set consuldnsbackup/data "$SERVICES_DATA"
  logger -t consul consuldnsbackup: Updated KV data.
fi
