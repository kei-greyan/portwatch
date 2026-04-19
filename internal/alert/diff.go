package alert

// Diff compares two sets of open ports and returns alerts for any changes.
// previous and current are slices of open port numbers.
func Diff(previous, current []int) []Alert {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var alerts []Alert

	// Ports that appeared in current but not in previous.
	for port := range currSet {
		if !prevSet[port] {
			alerts = append(alerts, PortOpened(port))
		}
	}

	// Ports that were in previous but missing from current.
	for port := range prevSet {
		if !currSet[port] {
			alerts = append(alerts, PortClosed(port))
		}
	}

	return alerts
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
