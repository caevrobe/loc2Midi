package main

import (
	"fmt"
	"loc2Midi/cfg"
	"loc2Midi/http"
	"loc2Midi/midi"
	"log"
	"math"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Bounds struct {
	center Coords // center of bounds
	offset Coords // not sure what this represents but I need it
	sw     Coords // bottom left boundary point
	ne     Coords // top right boundary point
}

type Coords struct {
	x float64
	y float64
}

var bounds Bounds
var absolute []Coords

func main() {
	stop := make(chan os.Signal, 1)
	recv := make(chan string)

	signal.Notify(stop, syscall.SIGINT)

	if err := cfg.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	writer, err := midi.GetWriter()
	if err != nil {
		log.Fatal(err)
	}

	initSpace()

	go http.ListenAndServe(cfg.Config.Port, recv)
	go liveCoords2Midi(recv, writer)

	//controller := exec.Command("cmd", "/C", "GRAPH.lnk") // todo prompt app install?
	controller := exec.Command("cmd", "/C", "explorer http://127.0.0.1:7777/")
	if err := controller.Start(); err != nil {
		log.Fatal(err)
	}

	<-stop
	writer.Close()
}

// initializes the listening space
func initSpace() {
	bounds.sw = Coords{0, 0} // todo replace with first GPS point
	bounds.ne = Coords{4, 4} // replace with second
	bounds.center = Coords{(bounds.ne.x + bounds.sw.x) / 2, (bounds.ne.y + bounds.sw.y) / 2}
	bounds.offset = Coords{(bounds.ne.x - bounds.sw.x) / 2, (bounds.ne.y - bounds.sw.y) / 2}

	// absolute sound coords (real space)
	absolute = make([]Coords, len(cfg.Config.TrackCoords))
	for i, tr := range cfg.Config.TrackCoords {
		absolute[i].x = tr[0]*bounds.offset.x + bounds.center.x
		absolute[i].y = tr[1]*bounds.offset.y + bounds.center.y
	}
}

func liveCoords2Midi(recv chan string, writer *midi.Writer) {
	var lat, lng, x, y, r_x, r_y, angle, distance float64
	var angle_send, distance_send int

	for msg := range recv {
		fmt.Sscanf(msg, "%f, %f", &lat, &lng)

		for n := range cfg.Config.TrackCoords {
			// sound locations (in -1, 1 space)
			x = (absolute[n].x - lat) / bounds.offset.x
			y = (absolute[n].y - lng) / bounds.offset.y

			// sound locations (centered around 0, 0)
			r_x = absolute[n].x + bounds.center.x - lat
			r_y = absolute[n].y + bounds.center.y - lng

			//16382
			// -8191, 8191
			// -180, 180

			angle = math.Atan2(r_x-bounds.center.x, r_y-bounds.center.y) * 180 / math.Pi
			angle_send = int(angle / 180 * 8191)

			distance = math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
			distance_send = int(distance * 16382) // max distance value?????
			distance_send2 := distance_send
			if distance_send2 > 0 {
				distance_send2 *= -1
			}
			distance_send -= 8191

			distance_send2 = int(math.Min(float64(distance_send2), 0)) / 25

			writer.Send(n*3, angle_send)
			writer.Send(n*3+1, distance_send)
			writer.Send(n*3+2, distance_send2)
		}
	}
}
