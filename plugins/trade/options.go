package trade

import (
	"fmt"
	"strconv"
	"strings"
)

type options struct {
	PriceSort  bool
	VolumeSort bool
	Source     string
	Dest       string
	Ammount    int
	Market     string
}

func parseOptions(args []string) (options, error) {
	var err error
	opt := options{}

	var c []string
	for _, v := range args {
		if strings.HasPrefix(v, "--market=") {
			opt.Market = strings.TrimPrefix(v, "--market=")
		} else if strings.HasPrefix(v, "--sort=") {
			switch strings.TrimPrefix(v, "--sort=") {
			case "price":
				opt.PriceSort = true
			case "volume":
				opt.VolumeSort = true
			}
		} else {
			c = append(c, v)
		}
	}
	if len(c) == 3 {
		if opt.Ammount, err = strconv.Atoi(args[0]); err != nil {
			return opt, fmt.Errorf("Wrong value of first argument : %v not a number.", c[0])
		}
		opt.Source = c[1]
		opt.Dest = c[2]
	} else if len(c) == 2 {
		opt.Ammount = 1
		opt.Source = c[0]
		opt.Dest = c[1]
	} else {
		return opt, fmt.Errorf("Not enough arguments.")
	}
	return opt, err
}
