package fixture

func helper(x int) int { return x * 2 }

var MaxRetries = 3

const DefaultName = "lattice"

func scale(v int) int {
	return helper(v) * 2
}
