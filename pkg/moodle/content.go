package moodle

type Content struct {
	Type           string
	FileName       string
	FilePath       string
	FileSize       int
	FileUrl        string
	TimeCreated    int64
	TimeModified   int64
	SortOrder      int
	Mimetype       string
	IsExternalFile bool
	UserId         int
	Author         string
	License        string
}

type ContentsInfo struct {
	FilesCount     int
	FilesSize      int
	LastModified   int
	MimeTypes      []string
	RepositoryType string
}

type Module struct {
	Id                  int
	Url                 string
	Name                string
	Instance            int
	Description         string
	Visible             int
	UserVisible         bool
	VisibleOnCoursePage int
	ModIcon             string
	ModName             string
	ModPlural           string
	Indent              int
	OnClick             string
	AfterLink           string
	CustomData          string
	NoViewLink          bool
	Completion          int
	Contents            []Content
	ContentsInfo        ContentsInfo
}

type Section struct {
	Id                  int
	Name                string
	Visible             int
	Summary             string
	SummaryFormat       int
	Section             int
	HiddenByNumSections int
	UserVisible         bool
	Modules             []Module
}
