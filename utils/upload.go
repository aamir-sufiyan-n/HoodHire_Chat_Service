
package utils

import (
    "context"
    "time"
    "os"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(file interface{}, filename string) (string, error) {
    cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
    if err != nil {
        return "", err
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    result, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
        Folder: "hoodhire/chat",
    })
    if err != nil {
        return "", err
    }
    return result.SecureURL, nil
}