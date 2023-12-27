package utils

import "os"

const (
	HelperImageEnv  = "HELPER_IMAGE"
	RsyslogImageEnv = "RSYSLOG_IMAGE"
)

func GetHelperImage() string {
	return os.Getenv(HelperImageEnv)
}

func GetRsyslogImage() string {
	return os.Getenv(RsyslogImageEnv)
}
