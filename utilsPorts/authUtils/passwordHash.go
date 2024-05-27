package auth_utils

type IPasswordHasher interface {
	GenerateHash(password string) (string, error)
	ComparePasswordhash(password string, hash string) error
}
