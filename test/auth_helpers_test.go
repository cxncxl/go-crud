package test

import (
	"testing"

	"github.com/cxcnxl/go-crud/internal/auth_helpers"
)

func TestRandomSalt(t *testing.T) {
    first := auth_helpers.GenerateRandomSalt();
    second := auth_helpers.GenerateRandomSalt();
    third  := auth_helpers.GenerateRandomSalt();

    valid := true;
    valid = valid && first != second;
    valid = valid && second != third;
    valid = valid && first != third;

    if !valid {
        t.Error("Salt not generated random every time")
    }
}

func TestHashing(t *testing.T) {
    testPass := "my_test_password";

    hash := auth_helpers.HashPassword(
        testPass,
        auth_helpers.GenerateRandomSalt(),
    );

    valid, err := auth_helpers.VerifyPass(testPass, hash);
    if err != nil {
        t.Errorf("Verify hash throwed an error %v\n", err);
    }
    if !valid {
        t.Error("Verify hash returned false, expected true");
    }
}
