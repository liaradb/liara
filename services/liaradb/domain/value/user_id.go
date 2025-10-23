package value

type UserID struct {
	baseString
}

func NewUserID(value string) UserID {
	return UserID{baseString(value)}
}
