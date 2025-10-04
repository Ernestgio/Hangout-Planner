package enums

type HangoutStatus string

const (
	StatusPlanning  HangoutStatus = "PLANNING"
	StatusConfirmed HangoutStatus = "CONFIRMED"
	StatusExecuted  HangoutStatus = "EXECUTED"
	StatusCancelled HangoutStatus = "CANCELLED"
)
