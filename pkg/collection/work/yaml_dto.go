package work

type YamlDto struct {
	Id          uint64   `yaml:"id"`
	Title       string   `yaml:"title"`
	Kind        string   `yaml:"kind"`
	Description string   `yaml:"description"`
	UserId      uint64   `yaml:"user_id"`
	UserName    string   `yaml:"user_name"`
	Restriction string   `yaml:"restriction"`
	Ai          string   `yaml:"ai"`
	Original    bool     `yaml:"original"`
	Pages       uint64   `yaml:"pages"`
	Views       uint64   `yaml:"views"`
	Bookmarks   uint64   `yaml:"bookmarks"`
	Likes       uint64   `yaml:"likes"`
	Comments    uint64   `yaml:"comments"`
	Uploaded    string   `yaml:"uploaded"`
	Downloaded  string   `yaml:"downloaded"`
	SeriesId    uint64   `yaml:"series_id"`
	SeriesTitle string   `yaml:"series_title"`
	SeriesOrder uint64   `yaml:"series_order"`
	Tags        []string `yaml:"tags"`
}

func (work *Work) YamlDto() *YamlDto {
	return &YamlDto{
		Id:          work.Id,
		Title:       work.Title,
		Kind:        work.Kind.String(),
		Description: work.Description,
		UserId:      work.UserId,
		UserName:    work.UserName,
		Restriction: work.Restriction.String(),
		Ai:          work.AiKind.String(),
		Original:    work.IsOriginal,
		Pages:       work.PageCount,
		Views:       work.ViewCount,
		Bookmarks:   work.BookmarkCount,
		Likes:       work.LikeCount,
		Comments:    work.CommentCount,
		Uploaded:    work.UploadTime.String(),
		Downloaded:  work.DownloadTime.String(),
		SeriesId:    work.SeriesId,
		SeriesTitle: work.SeriesTitle,
		SeriesOrder: work.SeriesOrder,
		Tags:        work.Tags,
	}
}
