package service

import (
	"strconv"
	"strings"
)

// GetPrefix принимает IP-адрес и маску сети в виде строк и возвращает префикс IP-адреса,
// который представляет собой результат применения маски к IP-адресу.
func GetPrefix(inputIP string, inputMask string) (string, error) {
	// Разделяем IP-адрес и маску сети на октеты (части).
	ip := strings.Split(inputIP, ".")
	mask := strings.Split(inputMask, ".")
	var prefix string

	// Проходим по каждому октету IP-адреса и маски, выполняя операцию "И" для получения префикса.
	for index, ipOct := range ip {
		intIPOct, err := strconv.Atoi(ipOct)
		if err != nil {
			return "", err
		}
		intMaskOct, err := strconv.Atoi(mask[index])
		if err != nil {
			return "", err
		}

		// Добавляем полученный результат в префикс. Разделяем октеты точками.
		prefix += strconv.Itoa(intIPOct & intMaskOct)
		if index != len(ip)-1 {
			prefix += "."
		}
	}
	return prefix, nil
}
