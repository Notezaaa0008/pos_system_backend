package storeDto

type GetUserStoreRequest struct {
	Page   int    `json:"page" binding:"gte=0"`
	Limit  int    `json:"limit" binding:"gte=0"` 
	Search string `json:"search"`           
}