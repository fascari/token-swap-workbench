package testdata

type (
	Payload struct {
		AccountID uint32 `json:"account_id" validate:"gt=0"`
		InToken   string `json:"in_token"   validate:"required"`
	}
)

func ValidPayload() Payload {
	return Payload{
		AccountID: 2,
		InToken:   "NEX",
	}
}

func PayloadWithZeroAccountID() Payload {
	return Payload{
		AccountID: 0,
		InToken:   "NEX",
	}
}

func PayloadWithEmptyInToken() Payload {
	return Payload{
		AccountID: 2,
		InToken:   "",
	}
}

func InvalidPayload() Payload {
	return Payload{
		AccountID: 0,
		InToken:   "",
	}
}
