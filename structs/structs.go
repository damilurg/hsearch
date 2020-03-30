package structs

const (
	KindOffer       = "offer"
	KindPhoto       = "photo"
	KindDescription = "description"
)

type (
	// Chat - all users and communicate with bot in chats. Chat can be group,
	// supergroup or private (type)
	Chat struct {
		Id       int64
		Username string
		Title    string // in private chats, this field save user full name
		Enable   bool
		Type     string
	}

	// Offer - posted on the site.
	Offer struct {
		Id         uint64
		Created    int64
		Url        string
		Topic      string
		FullPrice  string
		Price      int
		Currency   string // all currency is lower
		Phone      string
		Rooms      string // todo: convert to int
		Body       string
		Images     int
		ImagesList []string
	}

	// Answer - это ManyToMany для хранения реакции пользователя на
	// объявдение
	Answer struct {
		Created int64
		Chat    uint64
		Offer   uint64
		Like    bool
		Dislike bool
		Skip    uint64
	}

	// Feedback - структура для обратной связи в надежде получать баг репорты
	// а не угрозы что я бизнес чей-то сломал
	Feedback struct {
		Username string
		Chat     int64
		Body     string
	}
)
