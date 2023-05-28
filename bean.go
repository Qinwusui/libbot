package main

type UserInfo struct {
	UserName string `json:"qq"`
	Pswd     string `json:"code"`
	Token    string `json:"token"`
}
type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type GroupMsgData struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}
type TypedMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
type CodeUser struct {
	QQ   int64 `json:"qq"`
	Sent bool  `json:"sent"`
}

type RecallData struct {
	Data    map[string]interface{} `json:"data"`
	RetCode int                    `json:"retcode"`
	Status  string                 `json:"status"`
}
type SearchItem struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	ImageUrl []string `json:"imgUrl"`
	ImageDes []string `json:"imgDes"`
}
type MSG map[string]interface{}

type ForwardMsg MSG
type SauceNAOResp struct {
	Header  Header   `json:"header"`
	Results []Result `json:"results"`
}

type Header struct {
	UserID            string           `json:"user_id"`
	AccountType       string           `json:"account_type"`
	ShortLimit        string           `json:"short_limit"`
	LongLimit         string           `json:"long_limit"`
	LongRemaining     int              `json:"long_remaining"`
	ShortRemaining    int              `json:"short_remaining"`
	Status            int              `json:"status"`
	ResultsRequested  string           `json:"results_requested"`
	Index             map[string]Index `json:"index"`
	SearchDepth       string           `json:"search_depth"`
	MinimumSimilarity float64          `json:"minimum_similarity"`
	QueryImageDisplay string           `json:"query_image_display"`
	QueryImage        string           `json:"query_image"`
	ResultsReturned   int              `json:"results_returned"`
}

type Index struct {
	Status   int `json:"status"`
	ParentID int `json:"parent_id"`
	ID       int `json:"id"`
	Results  int `json:"results"`
}

type Result struct {
	Header InnerHeader `json:"header"`
	Data   InnerData   `json:"data"`
}

type InnerHeader struct {
	Similarity string `json:"similarity"`
	Thumbnail  string `json:"thumbnail"`
	IndexID    int    `json:"index_id"`
	IndexName  string `json:"index_name"`
	Dupes      int    `json:"dupes"`
	Hidden     int    `json:"hidden"`
}

type InnerData struct {
	ExtURLs     []string    `json:"ext_urls,omitempty"`
	Title       string      `json:"title,omitempty"`
	PixivID     int         `json:"pixiv_id,omitempty"`
	MemberName  string      `json:"member_name,omitempty"`
	MemberID    int         `json:"member_id,omitempty"`
	Published   string      `json:"published,omitempty"`
	Service     string      `json:"service,omitempty"`
	ServiceName string      `json:"service_name,omitempty"`
	ID          string      `json:"id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	UserName    string      `json:"user_name,omitempty"`
	Company     string      `json:"company,omitempty"`
	GetchuID    string      `json:"getchu_id,omitempty"`
	Source      string      `json:"source,omitempty"`
	Creator     interface{} `json:"creator,omitempty"`
	EngName     string      `json:"eng_name,omitempty"`
	JpName      string      `json:"jp_name,omitempty"`
	Path        string      `json:"path,omitempty"`
	CreatorName string      `json:"creator_name,omitempty"`
	AuthorName  string      `json:"author_name,omitempty"`
	AuthorURL   string      `json:"author_url,omitempty"`
}

type CFDetail struct {
	Status   string       `json:"status"`
	UserName string       `json:"userName,omitempty"`
	From     string       `json:"from,omitempty"`
	Result   []CFContests `json:"result"`
}

type CFContests struct {
	ID                  int64     `json:"id"`
	ContestId           int       `json:"contestId"`
	CreationTimeSeconds int64     `json:"creationTimeSeconds"`
	RelativeTimeSeconds int64     `json:"relativeTimeSeconds"`
	Problem             CFProblem `json:"problem"`
	Author              CFAuthor  `json:"author"`
	ProgrammingLanguage string    `json:"programmingLanguage"`
	Verdict             string    `json:"verdict"`
	Testset             string    `json:"testset"`
	PassedTestCount     int       `json:"passedTestCount"`
	TimeConsumedMillis  int64     `json:"timeConsumedMillis"`
	MemoryConsumedBytes int64     `json:"memoryConsumedBytes"`
}
type CFProblem struct {
	ContestId int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}
type CFAuthor struct {
	ConstestId       int                      `json:"contestId"`
	Members          []map[string]interface{} `json:"members"`
	ParticipantType  string                   `json:"participantType"`
	Ghost            bool                     `json:"ghost"`
	StartTimeSeconds int64                    `json:"startTimeSeconds"`
}
