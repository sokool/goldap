package sanitizer

import (
	"time"
	"strconv"
	"fmt"
	"encoding/base64"
	"gopkg.in/ldap.v2"
)

type (
	filter func(*ldap.EntryAttribute) error
)

var (
	definitions = map[string]filter{
		//"Ztou": ZuluToUnix,
		//"Ztot": ZuluToTimestamp,
		//"Ltou": LdapToUnix,
		"base64": Base64,
	}
)

func Base64(e *ldap.EntryAttribute) error {
	for i, bytes := range e.ByteValues {
		e.Values[i] = base64.StdEncoding.EncodeToString(bytes)
	}

	return nil
}

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

func Run(name string, field *ldap.EntryAttribute) error {

	if filter, ok := definitions[name]; ok {
		return filter(field)

	}

	return fmt.Errorf("Filter %s not available\n", name)
}