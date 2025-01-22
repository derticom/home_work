package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	Email string `json:"Email"` //nolint:tagliatelle // Tag name defined in JSON.
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := DomainStat{}
	decoder := json.NewDecoder(r)

	for {
		user := User{}

		if err := decoder.Decode(&user); err != nil {
			if errors.Is(err, io.EOF) {
				return domainStat, nil
			}
			return nil, err
		}

		if !strings.HasSuffix(user.Email, "."+domain) {
			continue
		}

		domainStat[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
	}
}
