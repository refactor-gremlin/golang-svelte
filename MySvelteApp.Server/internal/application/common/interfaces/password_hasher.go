package interfaces

// PasswordHasher abstracts hashing/verification logic so we can swap implementations easily.
type PasswordHasher interface {
	HashPassword(password string) (hash string, salt string, err error)
	VerifyPassword(password, hash, salt string) (bool, error)
}
