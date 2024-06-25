package folder

type FWS int

func (f FWS) Value() string {
	return " "
}

func (f FWS) Fold(limit int) string {
	return "\r\n "
}

func (f FWS) Priority() int {
	return int(f)
}
