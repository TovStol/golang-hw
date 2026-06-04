//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty input", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		invalidData := `{"Id":1,"Name":"Test","Email":"test@example.com"`
		_, err := GetDomainStat(bytes.NewBufferString(invalidData), "com")
		require.Error(t, err)
		require.Contains(t, err.Error(), "unmarshal error")
	})

	t.Run("case insensitive domain matching", func(t *testing.T) {
		caseData := `{"Id":1,"Name":"Test","Email":"user@EXAMPLE.COM"}
{"Id":2,"Name":"Test2","Email":"user2@example.com"}`
		result, err := GetDomainStat(bytes.NewBufferString(caseData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"example.com": 2,
		}, result)
	})

	t.Run("multiple same domains", func(t *testing.T) {
		multiData := `{"Id":1,"Name":"Test","Email":"user1@test.com"}
{"Id":2,"Name":"Test2","Email":"user2@test.com"}
{"Id":3,"Name":"Test3","Email":"user3@test.com"}`
		result, err := GetDomainStat(bytes.NewBufferString(multiData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"test.com": 3,
		}, result)
	})

	t.Run("mixed domains", func(t *testing.T) {
		mixedData := `{"Id":1,"Name":"Test","Email":"user1@test.com"}
{"Id":2,"Name":"Test2","Email":"user2@test.org"}
{"Id":3,"Name":"Test3","Email":"user3@test.net"}
{"Id":4,"Name":"Test4","Email":"user4@test.com"}`
		result, err := GetDomainStat(bytes.NewBufferString(mixedData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"test.com": 2,
		}, result)
	})

	t.Run("empty lines in input", func(t *testing.T) {
		emptyLinesData := `{"Id":1,"Name":"Test","Email":"user1@test.com"}

{"Id":2,"Name":"Test2","Email":"user2@test.com"}

`
		result, err := GetDomainStat(bytes.NewBufferString(emptyLinesData), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"test.com": 2,
		}, result)
	})
}
