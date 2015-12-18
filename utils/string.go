package utils

// StringInSlice checks if a string is present in an array of strings.
func StringInSlice(a string, l []string) bool {
	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

// IndexStringInSlice checks if a string is present in an array of strings.
// It will return the index of the item in the slice and true if such a string is present
// or -1 and false if the string wasn't found.
func IndexStringInSlice(a string, l []string) (int, bool) {
	for i, b := range l {
		if b == a {
			return i, true
		}
	}
	return -1, false
}

// RemoveString removes a string a from the slice l
func RemoveStringInSlice(a string, l []string) ([]string, bool) {
	i, in := IndexStringInSlice(a, l)

	if in {
		l = append(l[:i], l[i:]...)
	}
	return l, in
}
