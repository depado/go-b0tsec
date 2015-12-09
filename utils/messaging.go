package utils

const maxBufferSize = 450

// SplitMessage splits a long message into smaller ones so they can all be sent
// without the server dicthing data. Naive word-wrapping using space only.
func SplitMessage(m string) []string {
	var res []string
	rest := m
	if len(m) <= maxBufferSize {
		return append(res, m)
	}
	for len(rest) > maxBufferSize {
		for i := maxBufferSize; i > 0; i-- {
			if rest[i] == ' ' {
				res = append(res, rest[:i])
				rest = rest[i+1:]
				break
			}
		}
	}
	if len(rest) > 0 {
		res = append(res, rest)
	}
	return res
}
