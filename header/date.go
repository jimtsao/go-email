package header

import (
	"fmt"
	"time"
)

const time_RFC5322 = "Mon, 2 Jan 2006 15:04:05 -0700"

// Date represents the 'Date' header field
//
// Syntax:
//
//	date-time       =   [ day-of-week "," ] date time [CFWS]
//	day-of-week     =   ([FWS] day-name)
//	day-name        =   "Mon" / "Tue" / "Wed" / "Thu" /
//						"Fri" / "Sat" / "Sun"
//	date            =   day month year
//	day             =   ([FWS] 1*2DIGIT FWS)
//	month           =   "Jan" / "Feb" / "Mar" / "Apr" /
//						"May" / "Jun" / "Jul" / "Aug" /
//						"Sep" / "Oct" / "Nov" / "Dec"
//	year            =   (FWS 4*DIGIT FWS)
//	time            =   time-of-day zone
//	time-of-day     =   hour ":" minute [ ":" second ]
//	hour            =   2DIGIT
//	minute          =   2DIGIT
//	second          =   2DIGIT
//	zone            =   (FWS ( "+" / "-" ) 4DIGIT)s
type Date time.Time

func (d Date) Name() string {
	return "Date"
}

func (d Date) Validate() error {
	return nil
}

func (d Date) String() string {
	t := time.Time(d).Format(time_RFC5322)
	return fmt.Sprintf("%s: %s\r\n", d.Name(), t)
}
