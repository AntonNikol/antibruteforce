package service

import (
	"testing"

	"github.com/AntonNikol/anti-bruteforce/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetPrefix(t *testing.T) {
	cases := []struct {
		name    string
		network entity.IPNetwork
		expRes  struct {
			prefix string
			err    error
		}
	}{
		{
			name: "valid IP and mask",
			network: entity.IPNetwork{
				IP:   "192.168.1.1",
				Mask: "255.255.255.0",
			},
			expRes: struct {
				prefix string
				err    error
			}{prefix: "192.168.1.0", err: nil},
		},
		{
			name: "another valid IP and mask",
			network: entity.IPNetwork{
				IP:   "88.147.254.238",
				Mask: "255.255.255.240",
			},
			expRes: struct {
				prefix string
				err    error
			}{prefix: "88.147.254.224", err: nil},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			prefix, err := GetPrefix(testCase.network.IP, testCase.network.Mask)

			// Используем assert для сокращения кода и улучшения читаемости
			assert.Equal(t, testCase.expRes.err, err)
			assert.Equal(t, testCase.expRes.prefix, prefix)
		})
	}
}
