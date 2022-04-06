package tools

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/labstack/echo"
)

func PlaySongOneByOne(song uint32, e *echo.Echo) error {
	files, err := ioutil.ReadDir("music")

	if err != nil {
		log.Println(err)

	}

	if len(files) == 0 {
		log.Println("No stored music.")

		return nil
	}

	file, err := os.Open("music/" + files[song].Name())
	e.GET("/currentSong", func(c echo.Context) error {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return c.String(http.StatusOK, dir+"/music/"+files[song].Name())
	})

	if err != nil {
		log.Println(err)

	}

	defer file.Close()

	d, err := mp3.NewDecoder(file)

	if err != nil {

		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)

	if err != nil {

		return err
	}

	defer c.Close()

	p := c.NewPlayer()

	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {

		return err
	}

	return nil
}
