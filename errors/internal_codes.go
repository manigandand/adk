package errors

// Details struct contains the application logic specific internal error codes,
// readable message and a link to a support doc.
type Details struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

// ErrorCode is a enum type where all the platform error codes are defined.
type ErrorCode struct {
	code string
	link string
}

// All the defined Internal ErrorCodes
var (
	NoActiveSubscription = ErrorCode{
		code: "L0401",
		link: "https://support.gopherhut.com/docs/error-codes-and-what-they-mean#L0401",
	}
)

// Form forms the Details struct for the error code receiver.
func (e ErrorCode) Form(desc string) *Details {
	return &Details{
		Code:        e.code,
		Description: desc,
		Link:        e.link,
	}
}
