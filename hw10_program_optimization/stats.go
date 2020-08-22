package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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
	u, err := getUsersEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

func getUsersEmails(r io.Reader) (result []string, funcErr error) {
	var user User
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	bufReader := bufio.NewReader(r)
	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		}

		if err = json.Unmarshal(line, &user); err != nil {
			funcErr = err
			break
		}

		result = append(result, user.Email)
	}

	return
}

func countDomains(emails []string, domain string) (DomainStat, error) {
	re := regexp.MustCompile("\\." + domain)
	result := make(DomainStat)

	for _, email := range emails {
		matched := re.MatchString(email)
		if !matched {
			continue
		}

		domain := strings.Split(email, "@")[1]
		result[strings.ToLower(domain)]++
	}

	return result, nil
}
