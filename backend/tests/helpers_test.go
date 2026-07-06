package tests

const testUserID = 1

func fptr(v float64) *float64 { return &v }

func approx(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}
