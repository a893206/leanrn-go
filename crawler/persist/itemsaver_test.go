package persist

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"learngo/crawler/model"
	"testing"
)

func TestSave(t *testing.T) {
	expected := model.Profile{
		Name:      "披着羊皮的小红帽",
		City:      "北京",
		Age:       37,
		Education: "硕士",
		Marriage:  "未婚",
		Height:    162,
		Income:    "12001-20000元",
	}

	id, err := save(expected)

	if err != nil {
		panic(err)
	}

	// TODO: Try to start up elastic search
	// here using docker go client.
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	resp, err := client.Get().
		Index("dating_profile").
		Type("zhenai").
		Id(id).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	t.Logf("%s", resp.Source)

	var actual model.Profile
	err = json.Unmarshal(resp.Source, &actual)

	if err != nil {
		panic(err)
	}

	if actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}
