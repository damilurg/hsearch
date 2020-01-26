package structs

type (
	// User - telegram пользователь
	User struct {
		Id       uint64
		Username string
		Chat     int64
		Enable   bool
	}
)
