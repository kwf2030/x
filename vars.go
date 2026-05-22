package x

import (
	"errors"
	"runtime"
	"time"
)

const IsWindows = runtime.GOOS == "windows"

const (
	DateFormat1 = "2006-01-02"
	DateFormat2 = "2006/01/02"
	DateFormat3 = "2006_01_02"
	DateFormat4 = "2006.01.02"

	TimeFormat    = "15:04"
	TimeFormatSec = "15:04:05"
	TimeFormatMs  = "15:04:05.000"

	DateTimeFormat1 = "2006-01-02 15:04"
	DateTimeFormat2 = "2006/01/02 15:04"
	DateTimeFormat3 = "2006_01_02 15:04"
	DateTimeFormat4 = "2006.01.02 15:04"

	DateTimeFormatSec1 = "2006-01-02 15:04:05"
	DateTimeFormatSec2 = "2006/01/02 15:04:05"
	DateTimeFormatSec3 = "2006_01_02 15:04:05"
	DateTimeFormatSec4 = "2006.01.02 15:04:05"

	DateTimeFormatMs1 = "2006-01-02 15:04:05.000"
	DateTimeFormatMs2 = "2006/01/02 15:04:05.000"
	DateTimeFormatMs3 = "2006_01_02 15:04:05.000"
	DateTimeFormatMs4 = "2006.01.02 15:04:05.000"
)

var TimeZoneSH, _ = time.LoadLocation("Asia/Shanghai")

var (
	ErrInvalidArgs = errors.New("invalid arguments")
	ErrMissingArgs = errors.New("missing arguments")
)
