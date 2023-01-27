package model

type Health struct {
	Health   string `json:"health"`
	Database string `json:"database"`
}

const (
	// すべて正常に稼働中
	StatusGreen = "green"
	// 正常に稼働していないものがある
	StatusOrange = "orange"
	// 正常に稼働していない
	StatusRed = "red"
)
