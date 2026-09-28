package shortener

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num uint64) string {
	if num == 0 {
		return "000000"
	}

	s := ""

	for num > 0 {
		s = string(charset[num%62]) + s
		num /= 62
	}

	for len(s) < 6 {
		s = "0" + s
	}

	return s
}
