package id

const (
	IsGroup uint64 = 1 << 63
)

func GroupIDToSession(groupId int64) uint64 {
	return IsGroup | uint64(groupId)
}

func UserIDToSession(userId int64) uint64 {
	return uint64(userId)
}
