package main

type Course struct {
	Ret  int64  `json:"ret"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

type Data struct {
	KckbData []KckbDatum `json:"kckbData"`
}

type KckbDatum struct {
	ID      string `json:"id"`
	Xnxq    string `json:"xnxq"`
	Tid     string `json:"tid"`
	Type    int64  `json:"type"`
	Xs      int64  `json:"xs"`
	Rqxl    string `json:"rqxl"`
	Sfwc    int64  `json:"sfwc"`
	Jxbid   string `json:"jxbid"`
	Jxbmc   string `json:"jxbmc"`
	Zctype  string `json:"zctype"`
	Kcmc    string `json:"kcmc"`
	Kcbh    string `json:"kcbh"`
	Zc      string `json:"zc"`
	Zcstr   string `json:"zcstr"`
	Croommc string `json:"croommc"`
	Tmc     string `json:"tmc"`
	Xqid    string `json:"xqid"`
	Xqmc    string `json:"xqmc"`
	Xingqi  int64  `json:"xingqi"`
	Djc     int64  `json:"djc"`
	Flag    int64  `json:"flag"`
	Source  string `json:"source"`
	Pkid    string `json:"pkid"`
	Xq      string `json:"xq"`
	Bjdm    string `json:"bjdm"`
	Kcxz    string `json:"kcxz"`
	Xdxz    string `json:"xdxz"`
	Ksxs    string `json:"ksxs"`
}
