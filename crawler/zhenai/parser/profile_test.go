package parser

import (
	"io/ioutil"
	"learngo/crawler/engine"
	"learngo/crawler/model"
	"testing"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	result := ParseProfile(contents, "https://album.zhenai.com/u/104733447", "披着羊皮的小红帽")
	if len(result.Items) != 1 {
		t.Errorf("Items should contain 1 element; but was %v", result.Items)
	}

	actual := result.Items[0]

	expected := engine.Item{
		Url:  "https://album.zhenai.com/u/104733447",
		Type: "zhenai",
		Id:   "104733447",
		Payload: model.Profile{
			Name:      "披着羊皮的小红帽",
			City:      "北京",
			Age:       37,
			Education: "硕士",
			Marriage:  "未婚",
			Height:    162,
			Income:    "12001-20000元",
		},
	}

	if actual != expected {
		t.Errorf("expected %v, but was %v", expected, actual)
	}
}
