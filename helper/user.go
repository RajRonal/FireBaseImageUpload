package helper

func CreateImageUrl(imagePath string, bucket string) string {
	URL := "https://storage.cloud.google.com/" + bucket + "/" + imagePath
	return URL
}
