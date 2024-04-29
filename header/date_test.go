package header_test

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/jimtsao/go-email/header"
	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	d := header.Date(time.Now())
	assert.Equal(t, "Date", d.Name(), "name")
	assert.NoError(t, d.Validate(), "validate")

	sydney, _ := time.LoadLocation("Australia/Sydney")
	utc, _ := time.LoadLocation("UTC")
	midway, _ := time.LoadLocation("Pacific/Midway")
	kiri, _ := time.LoadLocation("Pacific/Kiritimati")
	cases := map[time.Time]string{
		time.Date(1990, time.April, 3, 5, 30, 15, 20, sydney):   "Date: Tue, 3 Apr 1990 05:30:15 +1000\r\n",
		time.Date(2000, time.January, 2, 12, 40, 20, 15, utc):   "Date: Sun, 2 Jan 2000 12:40:20 +0000\r\n",
		time.Date(2012, time.September, 10, 2, 4, 0, 0, midway): "Date: Mon, 10 Sep 2012 02:04:00 -1100\r\n",
		time.Date(2015, time.December, 31, 23, 59, 59, 0, kiri): "Date: Thu, 31 Dec 2015 23:59:59 +1400\r\n",
	}

	for input, want := range cases {
		got := header.Date(input).String()
		assert.Equal(t, want, got)
	}
}
