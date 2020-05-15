package model

// DeleteParam delete param , used for get deleted body
type DeleteParam struct {
	Key      int64  `json:"key,string" binding:"required"`
	DVersion string `json:"d_version" binding:"required"`
}
