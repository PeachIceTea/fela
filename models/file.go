package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/PeachIceTea/fela/conf"
)

// Declare errors
var (
	ErrNoAudioStream = errors.New("no audio stream")
)

type FileSeek interface {
	io.Reader
	io.Seeker
}

type File struct {
	ID int64 `db:"id" json:"id"`

	Name     string  `db:"name" json:"name"`
	Hash     string  `db:"hash" json:"hash"`
	Codec    string  `db:"codec" json:"codec"`
	Duration float64 `db:"duration" json:"duration"`

	Metadata *FileInfo `db:"metadata" json:"metadata"`

	Audiobook *int64 `db:"audiobook" json:"audiobook"`
	Chapter   *int   `db:"chapter" json:"chapter"`

	CreatedAt string `db:"created_at" json:"created_at"`
}

func (f *File) Insert(fs FileSeek, c *conf.Config) (err error) {
	f.Hash, err = Hash(fs)
	if err != nil {
		return
	}

	path := c.PathFromHash(f.Hash)
	err = Store(path, fs)
	if err != nil {
		return
	}

	f.Metadata, err = FileMetadata(path)
	if err != nil {
		os.Remove(path)
		return
	}

	f.Codec, err = f.Metadata.Codec()
	if err != nil {
		os.Remove(path)
		return
	}

	f.Duration, err = f.Metadata.Duration()
	if err != nil {
		os.Remove(path)
		return
	}

	metadata, err := f.Metadata.JSON()
	if err != nil {
		os.Remove(path)
		return
	}
	res, err := c.DB.Exec(c.TemplateString("file_insert"), f.Name, f.Hash, f.Codec, f.Duration, metadata)
	if err != nil {
		//FIXME: Duplicates are currently deleted
		os.Remove(path)
		return
	}

	f.ID, err = res.LastInsertId()
	return
}

func Hash(fs FileSeek) (s string, err error) {
	_, err = fs.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	sha := sha1.New()
	_, err = io.Copy(sha, fs)
	if err != nil {
		return
	}

	return hex.EncodeToString(sha.Sum(nil)), nil
}

func Store(path string, fs FileSeek) (err error) {
	_, err = fs.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	df, err := os.Create(path)
	if err != nil {
		return
	}
	defer df.Close()

	_, err = io.Copy(df, fs)
	return
}

// FileInfo - ffprobe output
type FileInfo struct {
	Streams []struct {
		Index          int    `json:"index"`
		CodecName      string `json:"codec_name"`
		CodecLongName  string `json:"codec_long_name"`
		Profile        string `json:"profile,omitempty"`
		CodecType      string `json:"codec_type"`
		CodecTimeBase  string `json:"codec_time_base,omitempty"`
		CodecTagString string `json:"codec_tag_string"`
		CodecTag       string `json:"codec_tag"`
		SampleFmt      string `json:"sample_fmt,omitempty"`
		SampleRate     string `json:"sample_rate,omitempty"`
		Channels       int    `json:"channels,omitempty"`
		ChannelLayout  string `json:"channel_layout,omitempty"`
		BitsPerSample  int    `json:"bits_per_sample,omitempty"`
		RFrameRate     string `json:"r_frame_rate"`
		AvgFrameRate   string `json:"avg_frame_rate"`
		TimeBase       string `json:"time_base"`
		StartPts       int    `json:"start_pts"`
		StartTime      string `json:"start_time"`
		DurationTs     int64  `json:"duration_ts"`
		Duration       string `json:"duration"`
		BitRate        string `json:"bit_rate,omitempty"`
		MaxBitRate     string `json:"max_bit_rate,omitempty"`
		NbFrames       string `json:"nb_frames,omitempty"`
		Disposition    struct {
			Default int `json:"default"`
			Dub     int `json:panic: sql: expected 5 arguments, got 4
			dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			Language    string `json:"language"`
			HandlerName string `json:"handler_name"`
		} `json:"tags,omitempty"`
		Width              int    `json:"width,omitempty"`
		Height             int    `json:"height,omitempty"`
		CodedWidth         int    `json:"coded_width,omitempty"`
		CodedHeight        int    `json:"coded_height,omitempty"`
		HasBFrames         int    `json:"has_b_frames,omitempty"`
		SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
		DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
		PixFmt             string `json:"pix_fmt,omitempty"`
		Level              int    `json:"level,omitempty"`
		ColorRange         string `json:"color_range,omitempty"`
		ColorSpace         string `json:"color_space,omitempty"`
		ChromaLocation     string `json:"chroma_location,omitempty"`
		Refs               int    `json:"refs,omitempty"`
		BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
	} `json:"streams"`
	Chapters []struct {
		ID        int    `json:"id"`
		TimeBase  string `json:"time_base"`
		Start     int    `json:"start"`
		StartTime string `json:"start_time"`
		End       int    `json:"end"`
		EndTime   string `json:"end_time"`
		Tags      struct {
			Title string `json:"title"`
		} `json:"tags"`
	} `json:"chapters"`
	Format struct {
		Filename       string `json:"filename"`
		NbStreams      int    `json:"nb_streams"`
		NbPrograms     int    `json:"nb_programs"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		StartTime      string `json:"start_time"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		ProbeScore     int    `json:"probe_score"`
		Tags           struct {
			MajorBrand       string `json:"major_brand"`
			MinorVersion     string `json:"minor_version"`
			CompatibleBrands string `json:"compatible_brands"`
			Title            string `json:"title"`
			Artist           string `json:"artist"`
			AlbumArtist      string `json:"album_artist"`
			Composer         string `json:"composer"`
			Album            string `json:"album"`
			Date             string `json:"date"`
			Encoder          string `json:"encoder"`
			Comment          string `json:"comment"`
			Genre            string `json:"genre"`
			Copyright        string `json:"copyright"`
		} `json:"tags"`
	} `json:"format"`
}

// Scan - Turns string into AudiobookInfo struct - used by database/sql
func (f *FileInfo) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var source []byte

	switch src.(type) {
	case string:
		{
			source = []byte(src.(string))
		}

	case []byte:
		{
			source = src.([]byte)
		}

	default:
		{
			return errors.New("incompatible type for FileInfo")
		}
	}

	return json.Unmarshal(source, f)
}

// JSON - Returns struct data as json string
func (f *FileInfo) JSON() (j string, err error) {
	b, err := json.Marshal(f)
	if err != nil {
		return
	}
	j = string(b)
	return
}

// Codec - Returns audio codec
func (f *FileInfo) Codec() (c string, err error) {
	for _, stream := range f.Streams {
		if stream.CodecType == "audio" {
			c = stream.CodecName
			return
		}
	}
	err = ErrNoAudioStream
	return
}

// Duration - Returns duration
func (f *FileInfo) Duration() (d float64, err error) {
	for _, stream := range f.Streams {
		if stream.CodecType == "audio" {
			d, err = strconv.ParseFloat(stream.Duration, 64)
			return
		}
	}

	err = ErrNoAudioStream
	return
}

func FileMetadata(path string) (info *FileInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", path, "-print_format", "json", "-v", "quiet", "-show_format", "-show_chapters", "-show_streams")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	info = &FileInfo{}
	if err = json.NewDecoder(stdout).Decode(&info); err != nil {
		return
	}

	err = cmd.Wait()
	return
}
