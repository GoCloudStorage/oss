package model

type UploadReq struct {
	ContentRange ContentRange
	TotalSize    int
	MD5          string
	Key          string
	ChunkNumber  int
}

type ContentRange struct {
	Start int
	End   int
	Total int
	start int64
}
