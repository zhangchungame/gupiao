package main

import (
	"fmt"
	"encoding/json"
	"os"
	"gopkg.in/olivere/elastic.v3"
	"log"
	"reflect"
)

type Shards struct {
	Total int
	Successful int
	Failed int
}

type source struct {
	User string
	Message string
}
type hits struct {
	Index string `json:"_index"`
	Type string `json:"_type"`
	Id string `json:"_id"`
	Source source `json:"_source"`
}
type hitsout struct {
	Total int
	Max_score float64
	Hits []hits `json:"hits"`
}
type Esback struct {
	Took int
	Timed_out bool
	Shards Shards `json:"_shards"`
	Hits hitsout

}

type Tweet struct {
	User string
	Message string
	Retweets int
}
func main() {
	//url:="http://localhost:9200/_search";
	//req,err:=http.NewRequest("POST",url,strings.NewReader("{\"query\":{\"match_all\":{}}}"))
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//client:=&http.Client{}
	//resp,err:=client.Do(req)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//body, err := ioutil.ReadAll(resp.Body)
	//str:=string(body)
	//var user map[string]interface{}
	//err=json.Unmarshal(body, &user)
	//if err!=nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(str)
	//esdata:=Esback{}
	//json.Unmarshal(body,&esdata)
	//fmt.Println(esdata)
	//fmt.Println(user["hits"]["hits"]["_source"]["User"])

	//singleInstance.GetEsInstance()


//	str:=`{"took":4,"timed_out":false,"_shards":{"total":10,"successful":10,"failed":0},"hitsout":{"total":5,"max_score":1.0,"hits":[
//{"_index":"gupiao","_type":"rikxian","_id":"AVlOaDT9N2hWhOFhINwl","_score":1.0,"_source":{"User":"zc","Message":"helloworld"}},
//{"_index":"gupiao","_type":"rikxian","_id":"AVlOcXm0N2hWhOFhINws","_score":1.0,"_source":{"User":"zc","Message":"helloworld"}},
// {"_index":"gupiao","_type":"rikxian","_id":"AVlOZbNXN2hWhOFhINwj","_score":1.0,"_source":{"User":"zc","Message":"helloworld"}},
// {"_index":"gupiao","_type":"rikxian","_id":"AVlOaHJZN2hWhOFhINwm","_score":1.0,"_source":{"User":"zc","Message":"helloworld"}},
// {"_index":"gupiao","_type":"rikxian","_id":"AVlOZ4-6N2hWhOFhINwk","_score":1.0,"_source":{"User":"zc","Message":"helloworld"}}]}}`
//	esdata:=Esback{}
//	err:=json.Unmarshal([]byte(str),&esdata)
//	if err!=nil {
//		fmt.Println(err)
//	}
//	fmt.Println(esdata)


	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("twitter").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("twitter").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Index a tweet (using JSON serialization)
	tweet1 := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("1").
		BodyJson(tweet1).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Index a second tweet (by string)
	tweet2 := `{"user" : "olivere", "message" : "It's a Raggy Waltz"}`
	put2, err := client.Index().
		Index("twitter").
		Type("tweet").
		Id("2").
		BodyString(tweet2).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put2.Id, put2.Index, put2.Type)

	// Get tweet with specified ID
	get1, err := client.Get().
		Index("twitter").
		Type("tweet").
		Id("1").
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("twitter").Do()
	if err != nil {
		panic(err)
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do()                // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		t := item.(Tweet)
		fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	script := elastic.NewScript("ctx._source.retweets += num").Param("num", 1)
	update, err := client.Update().Index("twitter").Type("tweet").Id("1").
		Script(script).
		Upsert(map[string]interface{}{"retweets": 0}).
		Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("New version of tweet %q is now %d", update.Id, update.Version)

	// ...

	// Delete an index.
	deleteIndex, err := client.DeleteIndex("twitter").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}