package util

import "hash/fnv"

// HashFnv calculates 32-bit fnv.New32a hash, given a string s
func HashFnv(s string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte (s))
	return hash.Sum32()
}
