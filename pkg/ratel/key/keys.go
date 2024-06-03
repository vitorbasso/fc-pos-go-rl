package key

type RatelKey string

func TokenKey(k string) RatelKey {
	return RatelKey("token:" + k)
}

func IPKey(k string) RatelKey {
	return RatelKey("ip:" + k)
}

func (k RatelKey) String() string {
	return string(k)
}
