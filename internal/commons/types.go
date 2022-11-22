package commons

type User string
type Offset uint
type Limit uint
type EventsCount uint

type Uid string
type FullName string
type Email string
type Start string
type End string
type Summary string
type GeoLat float32
type GeoLon float32

type Event struct {
	Uid      Uid
	FullName FullName
	Email    Email
	Start    Start
	End      End
	Summary  Summary
	GeoLat   GeoLat
	GeoLon   GeoLon
}

type GeneratorConfig struct {
	EventsCount               EventsCount            `json:"events_count"`
	ModificationRateInSeconds int                    `json:"modification_rate_in_seconds"`
	ReplacementRate           EventsCount            `json:"replacement_rate"`
	UpdateRate                EventsCount            `json:"update_rate"`
	Users                     []User                 `json:"users"`
	Options                   GeneratorConfigOptions `json:"options"`
}

type GeneratorConfigOptions struct {
	FullNames []FullName `json:"full_names"`
	Emails    []Email    `json:"emails"`
	Summaries []Summary  `json:"summaries"`
}
