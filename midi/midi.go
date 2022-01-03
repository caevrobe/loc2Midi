package midi

import (
	"errors"
	"loc2Midi/cfg"
	"log"
	"strings"

	writer "gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/midicatdrv"
)

type Writer struct {
	drv *driver.Driver
	*writer.Writer
}

// get loopMIDI writer
func GetWriter() (*Writer, error) {
	drv, err := driver.New()
	if err != nil {
		return nil, err
	}

	outs, err := drv.Outs()
	if err != nil {
		return nil, err
	}

	out := -1
	for x, port := range outs {
		if strings.Contains(port.String(), cfg.Config.MIDIPortName) {
			out = x
			break
		}
	}

	if out == -1 {
		return nil, errors.New("couldn't find MIDI output with name " + cfg.Config.MIDIPortName)
	}

	outs[out].Open() // todo add to struct, close when done

	return &Writer{drv, writer.New(outs[out])}, nil
}

func (wr *Writer) Close() {
	wr.drv.Close()
}

func (wr *Writer) Send(channel int, val int) {
	wr.SetChannel(uint8(channel))
	if err := writer.Pitchbend(wr, int16(val)); err != nil {
		log.Fatal(err)
	}
}
