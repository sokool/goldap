package sanitizer

import (
	"time"
	"strconv"
	"fmt"
	"encoding/base64"
	"gopkg.in/ldap.v2"
)

type (
	Sanitizer struct {
		fields map[string][]string
	}
)

var (
	filters = map[string]func(*ldap.EntryAttribute){
		//"Ztou": ZuluToUnix,
		//"Ztot": ZuluToTimestamp,
		//"Ltou": LdapToUnix,
		"base64encode": base64encode,
	}

	sanitizers = map[string]func(string) string{
		"base64decode": base64decode,
	}
)

func (s *Sanitizer) Register(sanitizer string, names []string) {
	for _, n := range names {
		s.fields[n] = append(s.fields[n], sanitizer)
	}
}

func (s *Sanitizer) Filter(name string, e *ldap.EntryAttribute) {
	if m, is := s.fields[name]; is {
		for _, op := range m {
			filters[op](e)
		}
	}
}

func (s *Sanitizer) Sanitize(name string, value string) (string, error) {
	if m, is := s.fields[name]; is {
		for _, op := range m {
			return sanitizers[op](value), nil
		}
	}

	return value, fmt.Errorf("Sanitizer %s not exists", s)
}

func base64encode(e *ldap.EntryAttribute) {
	for i, bytes := range e.ByteValues {
		e.Values[i] = base64.StdEncoding.EncodeToString(bytes)
	}
}

func base64decode(s string) string {
	o, _ := base64.StdEncoding.DecodeString(s)

	return string(o)
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

func New() *Sanitizer {
	return &Sanitizer{
		fields: make(map[string][]string),
	}
}