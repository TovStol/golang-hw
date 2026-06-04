package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var user User
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}

		if strings.HasSuffix(strings.ToLower(user.Email), "."+domain) {
			parts := strings.SplitN(user.Email, "@", 2)
			if len(parts) == 2 {
				fullDomain := strings.ToLower(parts[1])
				result[fullDomain]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return result, nil
}
