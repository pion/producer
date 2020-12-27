package ivf

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/ivfreader"
)

type IVFProducer struct {
	name    string
	stop    bool
	Samples chan media.Sample
	Track   *webrtc.TrackLocalStaticSample
	offset  int
}

func NewIVFProducer(name string, offset int) *IVFProducer {
	// Create track
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 90000}, "video", "pion")
	if err != nil {
		panic(err)
	}

	return &IVFProducer{
		name:    name,
		Samples: make(chan media.Sample),
		Track:   videoTrack,
		offset:  offset,
	}
}

func (t *IVFProducer) AudioTrack() *webrtc.TrackLocalStaticSample {
	return nil
}

func (t *IVFProducer) VideoTrack() *webrtc.TrackLocalStaticSample {
	return t.Track
}

func (t *IVFProducer) Stop() {
	t.stop = true
}

func (t *IVFProducer) SeekP(ts int) {
}

func (t *IVFProducer) Pause(pause bool) {
}

func (t *IVFProducer) Start() {
	go t.ReadLoop()
}

func (t *IVFProducer) VideoCodec() string {
	return webrtc.MimeTypeVP8
}

func (t *IVFProducer) ReadLoop() {
	startSeekFrames := t.offset * 30

	file, err := os.Open(t.name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ivf, header, ivfErr := ivfreader.NewWith(file)
	if ivfErr != nil {
		panic(ivfErr)
	}

	// Discard frames
	for i := 0; i < startSeekFrames; i++ {
		// TODO check for errors
		ivf.ParseNextFrame()
	}

	// Send our video file frame at a time. Pace our sending so we send it at the same speed it should be played back as.
	// This isn't required since the video is timestamped, but we will such much higher loss if we send all at once.
	sleepTime := time.Millisecond * time.Duration((float32(header.TimebaseNumerator)/float32(header.TimebaseDenominator))*1000)
	log.Println("Sleep time", sleepTime)
	for !t.stop {
		// Push sample
		frame, _, ivfErr := ivf.ParseNextFrame()
		if ivfErr == io.EOF {
			log.Println("All frames parsed and sent. Restart file")
			// TODO cleanup
			file.Seek(0, 0)
			ivf, header, ivfErr = ivfreader.NewWith(file)
			if ivfErr != nil {
				panic(ivfErr)
			}
			continue
		}

		if ivfErr != nil {
			log.Println("IVF error", ivfErr)
		}

		time.Sleep(sleepTime)
		if ivfErr = t.Track.WriteSample(media.Sample{Data: frame}); ivfErr != nil {
			log.Println("Track write error", ivfErr)
		}
	}
}
