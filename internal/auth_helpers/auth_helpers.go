package auth_helpers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
    "errors"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(data any) string {
    key := os.Getenv("JWT_SECRET");
    if key == "" {
        // TODO: verify this on earlier step
        panic("no jwt env variable set!");
    }

    dataMarshalled, err := json.Marshal(data);
    if err != nil {
        panic(err);
    }

    var dataMap map[string]any;
    err = json.Unmarshal(dataMarshalled, &dataMap);
    if err != nil {
        panic(err);
    }

    var payload jwt.MapClaims;
    for k := range dataMap {
        payload[k] = dataMap[k];
    }

    token := jwt.NewWithClaims(
        jwt.SigningMethodES256, 
        payload,
    );

    signed, err := token.SignedString([]byte(key));
    if err != nil {
        panic(err);
    }

    return signed;
}

func DecodeJWT(inp string) (map[string]any, error) {
    key := os.Getenv("JWT_SECRET");
    if key == "" {
        // TODO: verify this on earlier step
        panic("no jwt env variable set!");
    }

    var claims jwt.MapClaims;

    token, err := jwt.ParseWithClaims(
        inp,
        &claims,
        func(token *jwt.Token) (any, error) {
            return []byte(key), nil;
        },
    );
    if err != nil {
        if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
            return map[string]any{}, InvalidJwtError;
        }

        return map[string]any{}, err;
    }
    if !token.Valid {
        return map[string]any{}, InvalidJwtError;
    }

    return claims, nil;
}

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
        return false, UnknownSaltError;
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

var UnknownSaltError error;
var InvalidJwtError error;
