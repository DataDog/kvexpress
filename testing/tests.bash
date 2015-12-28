
trap "make consul_kill && make wercker_clean" INT TERM EXIT

export PREDICTED_CHECKSUM="0ab71a1c8fef24ade8d650e2cc248aac1e499a45a0e9456ba0b47901f99176d8"
export KVEXPRESS_DEBUG=1
export STOP_KEY_CONTENTS="Setting a stop key."
export LOCK_FILE=$(echo `pwd`/lock-test)
export LOCKED_FILE="$LOCK_FILE.locked"
export HOSTNAME=$(hostname)
export LOCK_URL="kvexpress/locks/e9de90b0a8985bf058580aa8457883b9fd76dd1fb13bdda5e1608884f9276dec/$HOSTNAME"
export URL_CHECKSUM="307b198c768b7a174b11e00c70bb1bd7b32597a86790279f763c4544dc12d1ff"
echo "This is a test of the lock-test file." > $LOCK_FILE

sleep 5

T_10insertSortingIntoTesting() {
  result="$(bin/kvexpress in -k testing -f sorting --sorted true)"
  [[ "$?" -eq "0" ]]
}

T_12testPredictedChecksum() {
  checksum="$(consul-cli kv-read kvexpress/testing/checksum)"
  [[ "$checksum" == "$PREDICTED_CHECKSUM" ]]
}

T_14createOutputFile() {
  bin/kvexpress out -k testing -f output
  checksum2="$(shasum -a 256 output | cut -d ' ' -f 1)"
  [[ "$checksum2" == "$PREDICTED_CHECKSUM" ]]
}

T_16createStopKey() {
  bin/kvexpress stop -k testing -r "$STOP_KEY_CONTENTS"
  stopkey="$(consul-cli kv-read kvexpress/testing/stop)"
  [[ "$stopkey" == "$STOP_KEY_CONTENTS" ]]
}

T_20tryToPullKey() {
  nostoppedfile="$(bin/kvexpress out -k testing -f stopped --verbose | grep 'Setting a stop key.')"
  [[ "$?" -eq "0" ]]
}

T_22testForStoppedFile() {
  [[ ! -f "stopped" ]]
}

T_24ignoreStopKey() {
  bin/kvexpress out -k testing -f ignored --ignore_stop
  [[ -f "ignored" ]]
}

T_30lockKey() {
  bin/kvexpress lock -f $LOCK_FILE
  lockkey="$(consul-cli kv-read $LOCK_URL | grep 'No reason given')"
  [[ "$?" -eq "0" ]]
}

T_32lockedFile() {
  [[ -f $LOCKED_FILE ]]
}

T_34lockedFileContents() {
  lockedFileContents="$(grep 'No reason given' $LOCKED_FILE)"
  [[ "$?" -eq "0" ]]
}

T_36testLockWorking() {
  bin/kvexpress out -k testing -f $LOCK_FILE --ignore_stop
  [[ `wc -l $LOCK_FILE | cut -d ' ' -f 8` == 1 ]]
}

T_38unlockFile() {
  bin/kvexpress unlock -f $LOCK_FILE
  [[ ! -f $LOCKED_FILE ]]
}

T_40testClean() {
  bin/kvexpress clean -f sorting
  [[ ! -f sorting ]]
  [[ ! -f sorting.compare ]]
  [[ ! -f sorting.last ]]
}

T_50testURLIn() {
  bin/kvexpress in -k url -u https://gist.githubusercontent.com/darron/9753b203b32667484105/raw/e66ea4c28c59e54aa8234d742368ccf93527dce5/gistfile1.txt
  urlchecksum=$(consul-cli kv-read kvexpress/url/checksum)
  [[ "$urlchecksum" == "$URL_CHECKSUM" ]]
}

T_52outputURL() {
  bin/kvexpress out -k url -f url
  urlfilechecksum="$(shasum -a 256 url | cut -d ' ' -f 1)"
  [[ "$urlfilechecksum" == "$URL_CHECKSUM" ]]
}

T_60getRawKey() {
  bin/kvexpress raw -k kvexpress/url/checksum -f raw_checksum -l 1
  rawchecksum="$(cat raw_checksum)"
  [[ "$rawchecksum" == "$URL_CHECKSUM" ]]
}

T_70copyKey() {
  bin/kvexpress copy --keyfrom url --keyto copied
  copiedchecksum=$(consul-cli kv-read kvexpress/copied/checksum)
  [[ "$copiedchecksum" == "$URL_CHECKSUM" ]]
}
