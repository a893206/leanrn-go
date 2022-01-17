package parser

import (
	"io/ioutil"
	"learngo/crawler/model"
	"testing"
)

func TestParseProfile(t *testing.T) {
	contents, err := ioutil.ReadFile("profile_test_data.html")
	if err != nil {
		panic(err)
	}
	result := ParseProfile(contents, "披着羊皮的小红帽")
	if len(result.Items) != 1 {
		t.Errorf("Items should contain 1 element; but was %v", result.Items)
	}

	profile := result.Items[0].(model.Profile)
	expected := model.Profile{
		Name:      "披着羊皮的小红帽",
		City:      "北京",
		Age:       37,
		Education: "硕士",
		Marriage:  "未婚",
		Height:    162,
		Income:    "12001-20000元",
	}
	if profile != expected {
		t.Errorf("expected %v, but was %v", expected, profile)
	}
}
