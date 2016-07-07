package sanitizer

import (
	"time"
	"strconv"
	"fmt"
)

type (
	sanitizer func(value, options []string) []string
)

var (
	helpers = map[string]sanitizer{
		"Ztou": ZuluToUnix,
		"Ltou": LdapToUnix,
	}
)

func ZuluToUnix(v, o []string) []string {
	zuluTime := v[0]

	location, _ := time.LoadLocation("Local")
	year, _ := strconv.Atoi(zuluTime[:4])
	month, _ := strconv.Atoi(zuluTime[4:6])
	day, _ := strconv.Atoi(zuluTime[6:8])
	hour, _ := strconv.Atoi(zuluTime[8:10])
	minute, _ := strconv.Atoi(zuluTime[10:12])
	second, _ := strconv.Atoi(zuluTime[12:14])

	time := time.Date(year, time.Month(month), day, hour, minute, second, 0, location).Format(time.UnixDate)

	return []string{time}
}

func LdapToUnix(v, o []string) []string {
	logonDate, _ := strconv.Atoi(v[0])
	time := time.Unix((int64(logonDate) / 10000000) - 11644473600, 0).Format(time.UnixDate);

	return []string{time}
}

func TimeToZulu(t time.Time) (string) {
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d.Z", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func Sanitize(name string, value, options []string) ([]string, bool) {

	if helper, ok := helpers[name]; ok {
		return helper(value, options), true
	}

	return nil, false
}