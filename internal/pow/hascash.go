package pow

import (
	"crypto/sha1"
	"fmt"
)

const zeroByte = 48

// HashCash - struct with data to find hashcash, implemented according to:
// https://en.wikipedia.org/wiki/Hashcash
type HashCash struct {
	Version    int
	ZerosCount int
	Date       int64
	Resource   string
	Rand       string
	Counter    int
}

func (h HashCash) header() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

func (h HashCash) hashHeader() string {
	return sha1Hash(h.header())
}

// sha1Hash - calculates sha1 hash from given string
func sha1Hash(data string) string {
	h := sha1.New()
	bs := h.Sum([]byte(data))
	return fmt.Sprintf("%x", bs)
}

func IsValidHashCash(in HashCash) bool {
	return isHashCorrect(in.hashHeader(), in.ZerosCount)
}

func CalcHashCash(in HashCash, maxIterations int) (HashCash, error) {
	for in.Counter <= maxIterations || maxIterations <= 0 {
		hash := in.hashHeader()
		//fmt.Println(header, hash)
		if isHashCorrect(hash, in.ZerosCount) {
			return in, nil
		}
		// if hash don't have needed count of leading zeros, we are increasing counter and try next hash
		// bump counter value and calculate the next hash if don't have enough leading zeros
		in.Counter++
	}
	return in, fmt.Errorf("max iterations exceeded")
}

// isHashCorrect - checks that hash has leading <zerosCount> zeros
func isHashCorrect(hash string, zerosCount int) bool {
	if zerosCount > len(hash) {
		return false
	}
	for _, ch := range hash[:zerosCount] {
		if ch != zeroByte {
			return false
		}
	}
	return true
}
