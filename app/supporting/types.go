package supporting

type SupportRequestStatus string

const (
	SupportRequestStatusNew    = SupportRequestStatus("new")
	SupportRequestStatusClosed = SupportRequestStatus("closed")
)

var availableSupportRequestStatuses = []SupportRequestStatus{
	SupportRequestStatusNew,
	SupportRequestStatusClosed,
}

func TryToConvertToValidSupportRequestStatus(statusStr string) SupportRequestStatus {
	for _, s := range availableSupportRequestStatuses {
		if s == SupportRequestStatus(statusStr) {
			return SupportRequestStatus(statusStr)
		}
	}
	return ""
}

type SupportRequest struct {
	ID         int64                `json:"id"`
	CreatedTs  int64                `json:"createdTs"`
	FirstName  string               `json:"firstName"`
	SecondName string               `json:"secondName"`
	LastName   string               `json:"lastName"`
	Phone      string               `json:"phone"`
	Message    string               `json:"Message"`
	Status     SupportRequestStatus `json:"status"`
	ClosedTs   int64                `json:"closedTs"`
}
