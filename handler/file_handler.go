package handler

import (
	"errors"
	"snackable/cache"
	"snackable/domain/file"
	"snackable/ext/snackable"

	"github.com/rs/zerolog/log"
)

var (
	ErrFileNotFound    = errors.New("file not found")
	ErrFileNotFinished = errors.New("file not finished")
)

type FileHandler interface {
	Handle(fileID string) (file.Model, error)
}

type fileHandler struct {
	snackableAPI snackable.Interface
	cache        cache.FileCache
}

func NewFileHandler(snackableAPI snackable.Interface, cache cache.FileCache) FileHandler {
	return fileHandler{
		snackableAPI: snackableAPI,
		cache:        cache,
	}
}

func (fh fileHandler) Handle(fileID string) (file.Model, error) {
	if v, err := fh.cache.Get(fileID); err == nil {
		log.Debug().Str("file_id", fileID).Msg("found from cache")
		return v, nil
	}

	f, err := fh.findByID(fileID)
	if err != nil && err != ErrFileNotFound {
		log.Error().Err(err).Str("file_id", fileID).Msg("failed to find snackable file by id")
		return file.Model{}, err
	}

	if err == ErrFileNotFound {
		return file.Model{}, err
	}

	if !f.IsFinished() {
		return file.Model{}, ErrFileNotFinished
	}

	details, err := fh.snackableAPI.Details(f.FileID)
	if err != nil {
		log.Error().Err(err).Str("file_id", fileID).Msg("failed to query file details")
		return file.Model{}, err
	}

	segments, err := fh.snackableAPI.Segments(f.FileID)
	if err != nil {
		log.Error().Err(err).Str("file_id", fileID).Msg("failed to query file segments")
		return file.Model{}, err
	}

	m := file.Model{
		FileID:           f.FileID,
		ProcessingStatus: f.ProcessingStatus,
		FileName:         details.FileName,
		MP3Path:          details.MP3Path,
		OriginalFilePath: details.OriginalFilePath,
		SeriesTitle:      details.SeriesTitle,
		Segments:         fh.remapSegments(segments),
	}

	if err := fh.cache.Set(fileID, m); err != nil {
		log.Error().Err(err).Msg("failed to save model to cache")
	}

	return m, nil
}

func (fh fileHandler) findByID(fileId string) (file.Model, error) {
	offset := 0

	// TODO: cache
	for {
		allResp, err := fh.snackableAPI.All(0, offset)
		if err != nil {
			return file.Model{}, nil
		}

		if len(allResp) == 0 {
			return file.Model{}, ErrFileNotFound
		}

		resp, found := allResp.FindByID(fileId)
		if found {
			return file.Model{
				FileID:           resp.FileID,
				ProcessingStatus: resp.ProcessingStatus,
			}, nil
		}
		offset += 5
	}
}

func (fh fileHandler) remapSegments(segs []snackable.SegmentResponse) []file.Segment {
	segments := make([]file.Segment, 0, len(segs))
	for _, v := range segs {
		remapped := file.Segment{
			FileSegmentID: v.FileSegmentID,
			FileID:        v.FileID,
			SegmentText:   v.SegmentText,
			StartTime:     v.StartTime,
			Endtime:       v.Endtime,
		}
		segments = append(segments, remapped)
	}
	return segments
}
