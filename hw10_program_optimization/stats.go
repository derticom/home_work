package hw10programoptimization

import (
	"encoding/json"
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
			if err == io.EOF {
				return domainStat, nil
			}
			return nil, err
		}

		if !strings.HasSuffix(user.Email, "."+domain) {
			continue
		}

		num := domainStat[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
		num++
		domainStat[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
	}
}
