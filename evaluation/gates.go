package evaluation

func Nand(a, b bool) bool {
	return !(a && b)
}

func Not(a bool) bool {
	return !a
}

func And(a, b bool) bool {
	return a && b
}

func Or(a, b bool) bool {
	return a || b
}

func Xor(a, b bool) bool {
	return a != b
}

func Mux(a, b, sel bool) bool {
	if sel {
		return a
	}
	return b
}

func Dmux(a, sel bool) (bool, bool) {
	if sel {
		return false, a
	}
	return a, false
}
