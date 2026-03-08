package value

type ClientVersion struct {
	baseString
}

func NewClientVersion(value string) ClientVersion {
	return ClientVersion{baseString(value)}
}
