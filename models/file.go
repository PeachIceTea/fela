package models

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"

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

	Name  string `db:"name" json:"name"`
	Hash  string `db:"hash" json:"hash"`
	Codec string `db:"codec" json:"codec"`

	Audiobook int64 `db:"audiobook" json:"audiobook"`
	Chapter   int   `db:"chapter" json:"chapter"`

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

	f.Codec, err = Codec(path)
	if err != nil {
		os.Remove(path)
		return
	}

	res, err := c.DB.Exec(c.TemplateString("file_insert"), f.Name, f.Hash, f.Codec)
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

// Codec - Returns audiocodec of given file
func Codec(path string) (c string, err error) {
	var data struct {
		Streams []struct {
			CodecName string `json:"codec_name"`
			CodecType string `json:"codec_type"`
		} `json:"streams"`
	}

	cmd := exec.Command("ffprobe", "-i", path, "-print_format", "json", "-v", "quiet", "-show_streams")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	err = json.NewDecoder(stdout).Decode(&data)
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	for _, stream := range data.Streams {
		if stream.CodecType == "audio" {
			c = stream.CodecName
			return
		}
	}

	err = ErrNoAudioStream
	return
}
