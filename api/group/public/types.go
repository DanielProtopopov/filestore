package public

type (
	UploadRS struct {
		Status string `json:"status"`
		Data   struct {
			Url string `json:"url"`
		} `json:"data"`
		Errors []string `json:"errors"`
	}
)
