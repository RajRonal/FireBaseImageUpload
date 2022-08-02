package handlers

import (
	"FireBaseEx/helper"
	"FireBaseEx/models"
	cloud "cloud.google.com/go/storage"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
)

func UploadImage(writer http.ResponseWriter, request *http.Request) {
	route := models.App{}
	var err error
	//ctx := context.Background()
	route.Ctx = context.Background()
	credentialsFile := option.WithCredentialsFile("serviceAccountKey.json")
	//fmt.Println(credentialsFile)
	app, err := firebase.NewApp(route.Ctx, nil, credentialsFile)
	if err != nil {
		log.Fatalln(err)
	}

	route.Client, err = app.Firestore(route.Ctx)
	if err != nil {
		log.Fatalln(err)
	}

	route.Storage, err = cloud.NewClient(route.Ctx, credentialsFile)
	if err != nil {
		log.Fatalln(err)
	}

	file, fileHeader, err := request.FormFile("image")
	err = request.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}

	defer file.Close()

	imagePath := fileHeader.Filename

	bucket := "image-a5e55.appspot.com"

	wc := route.Storage.Bucket(bucket).Object(imagePath).NewWriter(route.Ctx)
	_, err = io.Copy(wc, file)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return

	}
	if err := wc.Close(); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	imageUrl := helper.CreateImageUrl(imagePath, bucket)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	errs := json.NewEncoder(writer).Encode(imageUrl)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return

	}
}
