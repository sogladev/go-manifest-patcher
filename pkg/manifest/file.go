package manifest

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

func calculateHash(filePath string, hasher hash.Hash) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func CalculateHashSHA256(filePath string) (string, error) {
	return calculateHash(filePath, sha256.New())
}

func CalculateHashMD5(filePath string) (string, error) {
	return calculateHash(filePath, md5.New())
}
