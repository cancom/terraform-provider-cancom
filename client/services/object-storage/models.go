package client_object_storage

type Bucket struct {
	BucketName        string `json:"bucketName"`
	AvailabilityClass string `json:"availabilityClass"`
}

type User struct {
	Id              string                 `json:"id"`
	Username        string                 `json:"username"`
	Description     string                 `json:"description"`
	Permissions     map[string]interface{} `json:"permissions"`
	AccessKeyId     string                 `json:"accessKey"`
	SecretAccessKey string                 `json:"secretKey"`
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
