package structs

type (
	// UsersToOffers - это ManyToMany для хранения реакции пользователя на
	// объявдение
	UsersToOffers struct {
		UserId  uint64
		OfferId uint64
		Like    bool
		Dislike bool
		Skip    uint64
	}
)
