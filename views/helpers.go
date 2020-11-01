package views

import (
	"Website/db"
	"Website/jwt"
	"Website/settings"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

//AuthenticateContext func
func AuthenticateContext(c *gin.Context) (string, string, error) {
	tokenCookie, err := c.Cookie(settings.TokenCookieName)
	if err != nil {
		return "", "", err
	}

	return jwt.GetUsernameAndEmailFromToken(tokenCookie)
}

//GetUserFromContext returns user model from context
func GetUserFromContext(c *gin.Context) db.User {
	_, email, _ := AuthenticateContext(c)
	return db.GetUserByEmail(email)
}

func renderTemplate(c *gin.Context, name string, obj map[string]interface{}, title ...string) {
	var data gin.H

	if len(title) > 0 {
		data = gin.H{
			"title": title[0] + " | " + settings.SiteName,
		}
	} else {
		data = gin.H{
			"title": settings.SiteName,
		}
	}
	data["siteName"] = settings.SiteName

	username, _, err := AuthenticateContext(c)
	data["isAuthenticated"] = (err == nil)
	data["currentUsername"] = username

	for key, value := range obj {
		data[key] = value
	}
	c.HTML(http.StatusOK, name+".html", data)
}

func generateAvatarFileName(ext string) string {
	var genName = func() string {
		return uniuri.New() + ext
	}

	files, err := ioutil.ReadDir(settings.AvatarUploadsDir)
	if err != nil {
		log.Panic(err)
	}
	fileName := genName()

	for containsFile(files, fileName) {
		fileName = genName()
	}
	return fileName
}

func containsFile(files []os.FileInfo, fileName string) bool {
	for _, file := range files {
		if file.Name() == fileName {
			return true
		}
	}

	return false
}

func getContentTypeExt(contentType string) (string, bool) {
	for ext, t := range settings.AvatarWhitelistedContentTypes {
		if t == contentType {
			return ext, true
		}
	}

	return "", false
}

func resizeAndSave(fileHeader *multipart.FileHeader, dst string) error {
	imageFile, err := fileHeader.Open()
	if err != nil {
		return err
	}

	image, _, err := image.Decode(imageFile)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	resized := resize.Resize(
		settings.AvatarResizeWidth,
		settings.AvatarResizeHeight,
		image,
		resize.Bicubic)

	switch path.Ext(dst) {
	case ".jpg":
		err = jpeg.Encode(out, resized,
			&jpeg.Options{
				Quality: 10,
			})
	case ".png":
		err = png.Encode(out, resized)
	default:
		log.Panic("Invalid extension")
	}
	return err
}
