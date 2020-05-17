## File Producers

A set of file playback producers that write to `webrtc.Track`(s).

### Producers  

gst: gstreamer supported codecs - [rtwatch](https://github.com/pion/rtwatch)

ivf: vp8 playback from [play-from-disk](https://github.com/pion/webrtc/tree/master/examples/play-from-disk)

webm: direct file playback with vp8 or vp9 + ogg

### Usage

Producers support the following interface.

```
type IFileProducer interface {
	VideoTrack() *webrtc.Track
	VideoCodec() string
	AudioTrack() *webrtc.Track
	Stop()
	Start()
}

```
