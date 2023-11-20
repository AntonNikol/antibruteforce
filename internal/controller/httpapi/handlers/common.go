//nolint:gofumpt
package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
)

var validatePattern *regexp.Regexp
var errIPAlreadyExist = errors.New("this IP network already exist")

func init() {
	// Здесь мы инициализируем регулярное выражение, которое будет использоваться
	// для проверки корректности IP-адресов и масок.
	validatePattern = regexp.MustCompile(`(?m)^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`)
}

func initHeaders(writer http.ResponseWriter) {
	// Эта функция инициализирует заголовки HTTP-ответа для установки типа контента JSON.
	writer.Header().Set("Content-Type", "application/json")
}

func ValidateIP(network entity.IPNetwork) bool {
	// ValidateIP проверяет, является ли переданный IP-адрес и маска корректными.
	if !isCorrectIP(network.IP) {
		return false
	}
	if !isCorrectMask(network.Mask) {
		return false
	}
	return true
}

func isCorrectIP(ip string) bool {
	// isCorrectIP используется для проверки корректности IP-адреса с использованием регулярного выражения.
	return validatePattern.MatchString(ip)
}

func isCorrectMask(mask string) bool {
	// isCorrectMask используется для проверки корректности маски с использованием регулярного выражения.
	return validatePattern.MatchString(mask)
}

func ValidateRequest(request entity.Request) bool {
	// ValidateRequest проверяет корректность полей запроса, таких как логин, пароль и IP-адрес.
	if request.Login == "" || request.Password == "" {
		return false
	}
	if !isCorrectIP(request.IP) {
		return false
	}
	return true
}
