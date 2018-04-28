package domain

func markLabels(s string, offs []int) (int, error) {
	n := 0
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if start >= i {
				return 0, ErrEmptyLabel
			}

			if n >= len(offs) {
				return 0, ErrTooManyLabels
			}

			offs[n] = start
			n++

			start = i + 1
		}
	}

	if start < len(s) {
		if n >= len(offs) {
			return 0, ErrTooManyLabels
		}

		offs[n] = start
		n++
	}

	return n, nil
}

func getLabel(s string, out []byte) (int, error) {
	j := 1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' && i == len(s)-1 {
			break
		}

		if c >= 'a' && c <= 'z' {
			c &= 0xdf
		}

		if j >= len(out) {
			return 0, ErrLabelTooLong
		}

		out[j] = c
		j++
	}

	out[0] = byte(j - 1)

	return j, nil
}
