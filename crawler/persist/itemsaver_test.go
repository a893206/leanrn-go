package persist

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"learngo/crawler/engine"
	"learngo/crawler/model"
	"testing"
)

func TestSave(t *testing.T) {
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

	// TODO: Try to start up elastic search
	// here using docker go client.
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}

	// Save expected item
	err = save(expected)

	if err != nil {
		panic(err)
	}

	// Fetch saved item
	resp, err := client.Get().
		Index("dating_profile").
		Type(expected.Type).
		Id(expected.Id).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	t.Logf("%s", resp.Source)

	var actual engine.Item
	json.Unmarshal(resp.Source, &actual)

	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile

	// Verify result
	if actual != expected {
		t.Errorf("got %v, expected %v", actual, expected)
	}
}
