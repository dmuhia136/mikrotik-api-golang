package mikrotik

import (
	"crypto/tls"
	"fmt"

	"github.com/go-routeros/routeros"
)

func Connect(router Router) (*routeros.Client, error) {
	address := fmt.Sprintf("%s:%d", router.Host, router.Port)

	if router.UseTLS {
		return routeros.DialTLS(
			address,
			router.Username,
			router.Password,
			&tls.Config{InsecureSkipVerify: true},
		)
	}

	return routeros.Dial(address, router.Username, router.Password)
}
