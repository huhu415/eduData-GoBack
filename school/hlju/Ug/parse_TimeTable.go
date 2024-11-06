package hljuUg

import (
	"encoding/json"

	"eduData/repository"
)

type response struct {
	Code    int       `json:"code"`
	Msg     *string   `json:"msg"`
	MsgEn   *string   `json:"msg_en"`
	Content []content `json:"content"`
}

type content struct {
	Xn       string  `json:"xn"`
	Xq       string  `json:"xq"`
	Xn1      *string `json:"xn1"`
	Xq1      *string `json:"xq1"`
	Dj       string  `json:"dj"`
	Xj       string  `json:"xj"`
	Ks       string  `json:"ks"`
	Djms     string  `json:"djms"`
	Xjms     string  `json:"xjms"`
	Sxw      string  `json:"sxw"`
	Kssj     string  `json:"kssj"`
	Jssj     string  `json:"jssj"`
	Ksjc     *string `json:"ksjc"`
	Jsjc     *string `json:"jsjc"`
	Xnxqmc   *string `json:"xnxqmc"`
	XnxqmcEn *string `json:"xnxqmc_en"`
	Xqj      *string `json:"xqj"`
	Kskssj   *string `json:"kskssj"`
	Ksjssj   *string `json:"ksjssj"`
	Kssy     string  `json:"kssy"`
	Pksy     string  `json:"pksy"`
	Jysy     string  `json:"jysy"`
	Xssx     string  `json:"xssx"`
	Pylx     *string `json:"pylx"`
	Jglx     *string `json:"jglx"`
	Djz      *string `json:"djz"`
	Rq       *string `json:"rq"`
	Language *string `json:"language"`
}

func ParseTimeTableData(data *[]byte) ([]repository.Course, error) {
	var response response
	if err := json.Unmarshal(*data, &response); err != nil {
		return nil, err
	}
	var courses []repository.Course
	return courses, nil
}
