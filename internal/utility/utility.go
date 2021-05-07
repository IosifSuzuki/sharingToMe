package utility

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetTimeZoneOffsetFromGTW() int {
	_, offset :=  time.Now().Zone()
	return offset / 3600
}

func NumberGroup(number, digit uint64) string {
	var resultText string
	for i := uint64(0); number / 10 != 0; i++ {
		var appendText = ""
		if i != 0 && i % digit == 0 {
			appendText = " "
		}
		resultText = fmt.Sprintf("%d%s%s", number % 10, appendText, resultText)
		number /= 10
	}
	return resultText
}

func CreateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}