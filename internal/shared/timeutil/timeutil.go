package timeutil

import "time"

var saoPaulo *time.Location

func init() {
	var err error
	saoPaulo, err = time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic("failed to load America/Sao_Paulo timezone: " + err.Error())
	}
}

func Now() time.Time {
	return time.Now().In(saoPaulo)
}
