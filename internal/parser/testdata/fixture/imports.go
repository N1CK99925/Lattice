package fixture

import (
	"fmt"
	"strings"
)

func sanitize(s string) string {
	return strings.TrimSpace(s)
}

func logBoth(c Cache) {
	fmt.Println("cache size", c.size)
}
