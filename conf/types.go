package conf

// M - Shortcut for map
type M map[string]interface{}

// AudiobookInfo - Structure for ffprobe data
type AudiobookInfo struct {
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
			Default         int `json:"default"`
			Dub             int `json:"dub"`
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

// Book - Stores database book info
type Book struct {
	ID     int64  `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Author string `db:"author" json:"author"`

	Description *string `db:"description" json:"description,omitempty"`
	Metadata    *string `db:"metadata" json:"metadata,omitempty"`

	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt *string `db:"updated_at" json:"updated_at,omitempty"`
}
