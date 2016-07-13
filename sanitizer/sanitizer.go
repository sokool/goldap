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
		"Ztot": ZuluToTimestamp,
		"Ltou": LdapToUnix,

	}
)

func zuluToTime(zulu string) time.Time {

	location, _ := time.LoadLocation("Local")
	year, _ := strconv.Atoi(zulu[:4])
	month, _ := strconv.Atoi(zulu[4:6])
	day, _ := strconv.Atoi(zulu[6:8])
	hour, _ := strconv.Atoi(zulu[8:10])
	minute, _ := strconv.Atoi(zulu[10:12])
	second, _ := strconv.Atoi(zulu[12:14])

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, location)
}

func ZuluToTimestamp(v, o []string) []string {
	time := zuluToTime(v[0]).Unix()
	return []string{strconv.Itoa(int(time))}
}

func ZuluToUnix(v, o []string) []string {

	time := zuluToTime(v[0]).Format(time.UnixDate)

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