package gst

import (
	"log"

	"github.com/pion/webrtc/v3"
)

type GSTProducer struct {
	name       string
	audioTrack *webrtc.TrackLocalStaticSample
	videoTrack *webrtc.TrackLocalStaticSample
	pipeline   *Pipeline
	paused     bool
}

func NewGSTProducer(path string) *GSTProducer {
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264, ClockRate: 90000}, "video", "pion")
	if err != nil {
		log.Fatal(err)
	}

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2}, "audio", "pion")
	if err != nil {
		log.Fatal(err)
	}

	pipeline := CreatePipeline(path, audioTrack, videoTrack)

	return &GSTProducer{
		videoTrack: videoTrack,
		audioTrack: audioTrack,
		pipeline:   pipeline,
	}
}

func (t *GSTProducer) AudioTrack() *webrtc.TrackLocalStaticSample {
	return t.audioTrack
}

func (t *GSTProducer) VideoTrack() *webrtc.TrackLocalStaticSample {
	return t.videoTrack
}

func (t *GSTProducer) SeekP(ts int) {
	t.pipeline.SeekToTime(int64(ts))
}

func (t *GSTProducer) Pause(pause bool) {
	if pause {
		t.pipeline.Pause()
	} else {
		t.pipeline.Play()
	}
}

func (t *GSTProducer) Stop() {
}

func (t *GSTProducer) Start() {
	t.pipeline.Start()
}

func (t *GSTProducer) VideoCodec() string {
	return webrtc.MimeTypeH264
}
