package file

type Model struct {
	FileID           string `json:"fileId"`
	ProcessingStatus string `json:"processingStatus"`
	FileName         string `json:"fileName"`
	MP3Path          string `json:"mp3Path"`
	OriginalFilePath string `json:"originalFilePath"`
	SeriesTitle      string `json:"seriesTitle"`

	Segments []Segment `json:"segments"`
}

func (m Model) IsFinished() bool {
	return m.ProcessingStatus == "FINISHED"
}

type Segment struct {
	FileSegmentID int    `json:"fileSegmentId"`
	FileID        string `json:"fileId"`
	SegmentText   string `json:"segmentText"`
	StartTime     int    `json:"startTime"`
	Endtime       int    `json:"endTime"`
}
