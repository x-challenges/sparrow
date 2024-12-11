package instruments

// SliceOfAddresses
func SliceOfAddresses(iterator Iterator) []string {
	var res = []string{}

	for instrument := range iterator {
		res = append(res, instrument.Address)
	}

	return res
}
