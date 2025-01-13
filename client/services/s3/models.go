package client_s3

type Bucket struct {
	BucketName        string `json:"bucketName"`
	AvailabilityClass string `json:"availabilityClass"`
}

type User struct {
	Id          string                 `json:"id"`
	Username    string                 `json:"username"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
}

type UserCreateRequest struct {
	Username    string                 `json:"username"`
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
}

type UserUpdateRequest struct {
	Description string                 `json:"description"`
	Permissions map[string]interface{} `json:"permissions"`
}
