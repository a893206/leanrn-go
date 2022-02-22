package parser

import (
	"learngo/crawler/engine"
	"learngo/crawler/model"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`<div class="des f-cl"[^>]*>([^<]+)</div>`)
var idUrlRe = regexp.MustCompile(`https://album.zhenai.com/u/([\d]+)`)

func ParseProfile(contents []byte, url string, name string) engine.ParseResult {
	profile := model.Profile{}
	profile.Name = name

	list := extractString(contents, re)
	age, err := strconv.Atoi(list[1][:strings.Index(list[1], "岁")])
	if err == nil {
		profile.Age = age
	}
	height, err := strconv.Atoi(list[4][:strings.Index(list[4], "cm")])
	if err == nil {
		profile.Height = height
	}

	profile.City = list[0]
	profile.Education = list[2]
	profile.Marriage = list[3]
	profile.Income = list[5]

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      string(idUrlRe.FindSubmatch([]byte(url))[1]),
				Payload: profile,
			},
		},
	}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) []string {
	var res []string
	match := re.FindSubmatch(contents)
	for _, s := range strings.Split(string(match[1]), " | ") {
		res = append(res, s)
	}
	return res
}
