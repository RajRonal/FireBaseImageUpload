package handlers

import (
	"FireBaseEx/models"
	cloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadImage(writer http.ResponseWriter, request *http.Request) {
	client := models.App{}
	var err error
	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(os.Getenv("FIRE_KEY")))
	//fmt.Println(credentialsFile)
	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return
	}

	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		logrus.Error(err)
		return
	}

	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		logrus.Error(err)
		return
	}

	file, fileHeader, err := request.FormFile("image")
	err = request.ParseMultipartForm(10 << 20)
	if err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()
	imagePath := fileHeader.Filename + strconv.Itoa(int(time.Now().Unix()))
	bucket := "image-a5e55.appspot.com"
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)

	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bucketStorage.Close(); err != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	signedUrl := &cloud.SignedURLOptions{
		Scheme:  cloud.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}
	url, err := client.Storage.Bucket(bucket).SignedURL(imagePath, signedUrl)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Println(url)
	errs := json.NewEncoder(writer).Encode(url)
	if errs != nil {
		logrus.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
