/*
Samuel Worne <samuel@worne.xyz>

sources
https://github.com/BurntSushi/xgbutil/blob/master/_examples/screenshot/main.go
https://gist.github.com/DavidVaini/10308388
http://stackoverflow.com/questions/976026/low-hanging-graphics-programming-fruits/976560#976560
*/

package main
import (
	"log"
	"os"
	//"fmt"
	"flag"
	"os/user"
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/draw"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func main () {

	var pixel_size = flag.Int("s", 16, "size of pixelation for blur")
	var filename = flag.String("f", saveDir(), "Overide default save directory")
	flag.Parse()
	//fmt.Println(getUserHostname())
	//fmt.Println(*quiet)
	//fmt.Println(*pixel_size)
	//fmt.Println(username.Username)
	//fmt.Println(current_user.Username + "@" + hostname)
	//var quiet = flag.Bool("q", true, "hide locked icon")

	ximg := captureScreen()
	img := pixelateImage(ximg, *pixel_size)
	saveImg(img, *filename)

}

func saveDir () (save_dir string ) {
	var buffer bytes.Buffer
	buffer.WriteString(getHomeDir())
	buffer.WriteString("/.lock.jpg")
	save_dir = buffer.String()
	return
}

func getUserHostname () (user_and_hostname string) {
	hostname, err := os.Hostname()
	current_user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
        buffer.WriteString(hostname)
        buffer.WriteString("@")
        buffer.WriteString(current_user.Username)
        user_and_hostname = buffer.String()
	return
}


// XSERVER
func captureScreen() (ximg image.Image) {
	X, err := xgbutil.NewConn()
	if err != nil {
	log.Fatal(err)
	}

	ximg, err = xgraphics.NewDrawable(X, xproto.Drawable(X.RootWin()))
	if err != nil {
		log.Fatal(err)
	}
	return
}


// IMAGE FUNCTIONS
func saveImg(img image.Image, filename string) {
	//var filename = getHomeDir() + "/.lock.jpg"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
    }
	defer file.Close()
	var opt jpeg.Options
        opt.Quality = 85
	jpeg.Encode(file, img, &opt)
}

func pixelateImage (img image.Image, pixel_size int) (pixelated_img image.Image) {
	image_x, image_y := getImageDimensions(img)
	pixelated_img = image.NewRGBA(image.Rect(0, 0, image_x, image_y))
	for count_x := 0; count_x < image_x; count_x = count_x + pixel_size {
		for  count_y := 0; count_y < image_y; count_y = count_y + pixel_size {
			pixel_color := getPixelColor(img, count_x, count_y, pixel_size, pixel_size)
			for i := count_y; i < count_y + pixel_size; i++ {
				for j := 0; j < pixel_size; j++ {
					pixelated_img.(draw.Image).Set(count_x + j, i, pixel_color)
				}
			}
		}
	}
	return
}

func getPixelColor (img image.Image, x, y, h, w int) (pixel_color color.Color) {
	var w_rounded int = (w / 2)
	var h_rounded int = (h / 2)
	var color_x int = x + w_rounded
	var color_y int = y + h_rounded
	pixel_color = img.At(color_x, color_y) //color.Color
	return
}

func getImageDimensions (img image.Image) (image_x, image_y int) {
	b := img.Bounds()
        image_x = b.Max.X
        image_y = b.Max.Y
	return
}

func getHomeDir() (home_dir string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	home_dir = usr.HomeDir
	return
}
