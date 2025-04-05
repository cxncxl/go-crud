package auth_helpers

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"strings"
)

func HashPassword(pass string, salt string) string {
    passSalted := fmt.Sprintf("%s.%s", salt, pass);
    passHashedBytes := md5.Sum([]byte(passSalted));

    var passHashed string;
    for _, b := range passHashedBytes {
        passHashed += string(rune(b));
    }

    hash := fmt.Sprintf("%s%s", salt, hashSaltDelimeter);
    hash = fmt.Sprintf("%s%s", hash, passHashed);

    return hash;
}

func VerifyPass(pass string, hash string) (bool, error) {
    parts := strings.Split(hash, hashSaltDelimeter);
    if len(parts) < 2 {
        return false, UnknownSaltError{};
    }

    salt := parts[0];

    validHash := HashPassword(pass, salt);

    return hash == validHash, nil;
}

func GenerateRandomSalt() string {
	salt := make([]rune, saltSize);

	for i := range salt {
		salt[i] = alphabet[rand.Intn(len(alphabet))];
	}

    return string(salt);
}

const saltSize int = 16;
const hashSaltDelimeter string = ";";

var alphabet = []rune(
    "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890!@#$%^&*()",
);

type UnknownSaltError struct {}
func (self UnknownSaltError) Error() string {
    return "Unknwon Salt / Invalid password format";
}
