package utils

const RsyslogConfigFormat = `$ModLoad imuxsock
$SystemLogSocketName %s
$ModLoad omstdout.so
*.* :omstdout:`
