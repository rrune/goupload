package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	// Convert password string to byte slice

	var passwordBytes = []byte(password)

	// Hash password with bcrypt's min cost

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func DoPasswordsMatch(hashedPassword, currPassword string) bool {

	err := bcrypt.CompareHashAndPassword(

		[]byte(hashedPassword), []byte(currPassword))

	return err == nil

}
