package commons

import (
	"encoding/json"
	"fmt"

	"strconv"

	"reflect"

	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

var (
	esHelper   *es
	master     *elastic.Client
	slave      *elastic.Client
	ctx        = context.Background()
	scrollKeep string
)

type es struct{}

type RangeQuery struct {
	Field string
	Gte   interface{}
	Gt    interface{}
	Lte   interface{}
	Lt    interface{}
}

func GetES() *es {
	if esHelper == nil {
		esHelper = &es{}
	}
	masterAddress := Config.GetString("es.master")
	slaveAddress := Config.GetString("es.slave")
	scrollKeep = Config.GetString("es.scrollKeep")
	master = initClient(masterAddress)
	slave = initClient(slaveAddress)
	return esHelper
}

func initClient(address string) *elastic.Client {
	var (
		client *elastic.Client
		err    error
	)
	client, err = elastic.NewClient(
		elastic.SetURL(address),
		elastic.SetSniff(false),
	)
	if err != nil {
		GetLogger().LogErr(err)
	}
	_, _, err = client.Ping(address).Do(ctx)
	if err != nil {
		GetLogger().LogErr(err)
	}
	return client
}

// func checkIndex(index string) error {
// 	exists, err := client.IndexExists(index).Do(context.Background())
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf(ES_INDEX_NOT_EXISTS)
// 	}
// 	return nil
// }

// @Title SearchAll
// @Description Search data from es
// @Parameters
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string                       size
//            sort             map[string]bool              sort map
// @Returns result:[]byte err:error
func (*es) SearchAll(index, _type, sizeStr string, sort map[string]bool) ([]byte, error) {
	var (
		list       []interface{}
		ess        *elastic.SearchService
		result     *elastic.SearchResult
		err        error
		object     interface{}
		byteResult []byte
		size       int
	)
	ess = slave.Search(index).Type(_type).Query(elastic.NewMatchAllQuery())
	if sort != nil {
		for k, v := range sort {
			if v {
				ess = ess.SortBy(elastic.NewFieldSort(k).Asc())
			} else {
				ess = ess.SortBy(elastic.NewFieldSort(k).Desc())
			}
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err == nil {
			ess = ess.Size(size)
		}
	}
	result, err = ess.Do(ctx)
	if err != nil {
		return nil, err
	}
	list = []interface{}{}
	if result != nil && result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, err
			}
			list = append(list, object)
		}
	} else {
		return nil, nil
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return byteResult, nil
}

// @Title FuzzyScroll
// @Description fuzzy search by scroll
// @Parameters
//            must             map[string]interface{}       must map
//            mustNot          map[string]interface{}       mustNot map
//            should           map[string]interface{}       should map
//            sort             map[string]bool              sort map
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string						page size
//            scrollID         string                       scroll id
// @Returns result:[]byte scrollId:string err:error
func (*es) FuzzyScroll(must, mustNot, should map[string]interface{}, sort map[string]bool, ranges []*RangeQuery, index, _type, sizeStr, scrollID string) ([]byte, string, error) {
	var (
		byteResult []byte
		ess        *elastic.ScrollService
		result     *elastic.SearchResult
		err        error
		object     interface{}
		list       []interface{}
		rq         *elastic.RangeQuery
		bq         *elastic.BoolQuery
		qlist      []elastic.Query
		size       int
	)
	ess = slave.Scroll(index).Type(_type).Scroll(scrollKeep)
	if scrollID == "" {
		bq = elastic.NewBoolQuery()
		if must != nil {
			flag := false
			qlist = []elastic.Query{}
			for k, v := range must {
				if flag {
					qlist = append(qlist, elastic.NewMatchQuery(k, v).Fuzziness("AUTO").Operator("AND"))
					continue
				}
				qlist = append(qlist, elastic.NewMatchQuery(k, v))
			}
			bq = bq.Must(qlist...)
		}
		if mustNot != nil {
			flag := false
			qlist = []elastic.Query{}
			for k, v := range mustNot {
				if flag {
					qlist = append(qlist, elastic.NewMatchQuery(k, v).Fuzziness("AUTO").Operator("AND"))
					continue
				}
				qlist = append(qlist, elastic.NewMatchQuery(k, v))
			}
			bq = bq.MustNot(qlist...)
		}
		if should != nil {
			flag := false
			qlist = []elastic.Query{}
			for k, v := range should {
				if flag {
					qlist = append(qlist, elastic.NewMatchQuery(k, v).Fuzziness("AUTO").Operator("AND"))
					continue
				}
				qlist = append(qlist, elastic.NewMatchQuery(k, v))
			}
			bq = bq.Should(qlist...)
		}
		if ranges != nil {
			for _, v := range ranges {
				rq = elastic.NewRangeQuery(v.Field)
				if v.Gte != nil {
					rq = rq.Gte(v.Gte)
				}
				if v.Gt != nil {
					rq = rq.Gt(v.Gt)
				}
				if v.Lte != nil {
					rq = rq.Lte(v.Lte)
				}
				if v.Lt != nil {
					rq = rq.Lt(v.Lt)
				}
				bq = bq.Must(rq)
			}
		}
		ess = ess.Query(bq)
		if sort != nil {
			for k, v := range sort {
				if v {
					ess = ess.SortBy(elastic.NewFieldSort(k).Asc())
				} else {
					ess = ess.SortBy(elastic.NewFieldSort(k).Desc())
				}
			}
		}
		if sizeStr != "" {
			size, err = strconv.Atoi(sizeStr)
			if err != nil {
				return nil, "", err
			}
			ess = ess.Size(size)
		}
		result, err = ess.Do(ctx)
	} else {
		result, err = ess.ScrollId(scrollID).Do(ctx)
	}
	if err != nil {
		return nil, "", err
	}
	list = []interface{}{}
	if result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, "", err
			}
			list = append(list, object)
		}
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, "", err
	}
	return byteResult, result.ScrollId, nil
}

// @Title FuzzySearch
// @Description fuzzy search
// @Parameters
//            must             map[string]interface{}       must map
//            mustNot          map[string]interface{}       mustNot map
//            should           map[string]interface{}       should map
//            ranges           []*RangeQuery                range slice
//            sort             map[string]bool              sort map
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string						page size
// @Returns result:[]byte err:error
func (*es) FuzzySearch(must, mustNot, should map[string]interface{}, ranges []*RangeQuery, sort map[string]bool, index, _type, sizeStr string) ([]byte, error) {
	var (
		byteResult []byte
		ess        *elastic.SearchService
		result     *elastic.SearchResult
		err        error
		object     interface{}
		list       []interface{}
		qlist      []elastic.Query
		rq         *elastic.RangeQuery
		bq         *elastic.BoolQuery
		size       int
	)
	ess = slave.Search(index).Type(_type)
	bq = elastic.NewBoolQuery()
	if must != nil {
		qlist = []elastic.Query{}
		for k, v := range must {
			qlist = append(qlist, elastic.NewMatchPhrasePrefixQuery(k, v))
		}
		bq = bq.Must(qlist...)
	}
	if mustNot != nil {
		qlist = []elastic.Query{}
		for k, v := range mustNot {
			qlist = append(qlist, elastic.NewMatchPhrasePrefixQuery(k, v))
		}
		bq = bq.MustNot(qlist...)
	}
	if should != nil {
		qlist = []elastic.Query{}
		for k, v := range should {
			qlist = append(qlist, elastic.NewMatchPhrasePrefixQuery(k, v))
		}
		bq = bq.Should(qlist...)
	}
	if ranges != nil {
		for _, v := range ranges {
			rq = elastic.NewRangeQuery(v.Field)
			if v.Gte != nil {
				rq = rq.Gte(v.Gte)
			}
			if v.Gt != nil {
				rq = rq.Gt(v.Gt)
			}
			if v.Lte != nil {
				rq = rq.Lte(v.Lte)
			}
			if v.Lt != nil {
				rq = rq.Lt(v.Lt)
			}
			bq = bq.Must(rq)
		}
	}
	ess = ess.Query(bq)
	if sort != nil {
		for k, v := range sort {
			if v {
				ess = ess.SortBy(elastic.NewFieldSort(k).Asc())
			} else {
				ess = ess.SortBy(elastic.NewFieldSort(k).Desc())
			}
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nil, err
		}
		ess = ess.Size(size)
	}
	result, err = ess.Do(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	list = []interface{}{}
	if result != nil && result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, err
			}
			list = append(list, object)
		}
	} else {
		return nil, nil
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return byteResult, nil
}

// @Title Search
// @Description search
// @Parameters
//            must             map[string]interface{}       must map
//            mustNot          map[string]interface{}       mustNot map
//            should           map[string]interface{}       should map
//            sort             map[string]bool              sort map
//            ranges           []*RangeQuery                range slice
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string						page size
// @Returns result:[]byte scrollId:string err:error
func (*es) Search(must, mustNot, should map[string]interface{}, sort map[string]bool, ranges []*RangeQuery, index, _type, sizeStr string) ([]byte, error) {
	var (
		ess        *elastic.SearchService
		result     *elastic.SearchResult
		err        error
		bq         *elastic.BoolQuery
		slist      []elastic.Query
		mlist      []elastic.Query
		nmlist     []elastic.Query
		list       []interface{}
		object     interface{}
		byteResult []byte
		rq         *elastic.RangeQuery
		size       int
	)
	ess = slave.Search(index).Type(_type)
	bq = elastic.NewBoolQuery()
	if must != nil {
		mlist = []elastic.Query{}
		for k, v := range must {
			mlist = append(mlist, elastic.NewMatchQuery(k, v))
		}
		bq = bq.Must(mlist...)
	}
	if should != nil {
		slist = []elastic.Query{}
		for k, v := range should {
			slist = append(slist, elastic.NewMatchQuery(k, v))
		}
		bq = bq.Should(slist...)
	}
	if mustNot != nil {
		nmlist = []elastic.Query{}
		for k, v := range mustNot {
			nmlist = append(nmlist, elastic.NewMatchQuery(k, v))
		}
		bq = bq.MustNot(nmlist...)
	}
	if sort != nil {
		for k, v := range sort {
			if v {
				ess = ess.SortBy(elastic.NewFieldSort(k).Asc())
			} else {
				ess = ess.SortBy(elastic.NewFieldSort(k).Desc())
			}
		}
	}
	if ranges != nil {
		for _, v := range ranges {
			rq = elastic.NewRangeQuery(v.Field)
			if v.Gte != nil {
				rq = rq.Gte(v.Gte)
			}
			if v.Gt != nil {
				rq = rq.Gt(v.Gt)
			}
			if v.Lte != nil {
				rq = rq.Lte(v.Lte)
			}
			if v.Lt != nil {
				rq = rq.Lt(v.Lt)
			}
			bq = bq.Must(rq)
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return nil, err
		}
		ess = ess.Size(size)
	}
	result, err = ess.Query(bq).Do(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	list = []interface{}{}
	if result != nil && result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, err
			}
			list = append(list, object)
		}
	} else {
		return nil, nil
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return byteResult, nil
}

// @Title SuggestSearch
// @Description Suggest Search
// @Parameters
//            suggest          string                       suggest key word
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string                       size
// @Returns result:[]byte err:error
func (*es) SuggestSearch(suggest string, index, _type string, sizeStr string) ([]byte, error) {
	var (
		list       []interface{}
		ess        *elastic.SearchService
		result     *elastic.SearchResult
		err        error
		object     interface{}
		byteResult []byte
		size       int
	)
	ess = slave.Search(index).Type(_type).Suggester(elastic.NewCompletionSuggester(suggest))
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err == nil {
			ess = ess.Size(size)
		}
	}
	result, err = ess.Do(ctx)
	if err != nil {
		return nil, err
	}
	list = []interface{}{}
	if result != nil && result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, err
			}
			list = append(list, object)
		}
	} else {
		return nil, nil
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return byteResult, nil
}

// @Title AssociativeSearch
// @Description Associative Search
// @Parameters
//            regexp           map[string]string            regexp map
//            postFilter       map[string]interface{}       postFilter map
//            sort             map[string]bool              sort map
//            index            string                       es index
//            _type            string                       es type
//            sizeStr          string                       size
// @Returns result:[]byte err:error
func (*es) AssociativeSearch(regexp map[string]string, postFilter map[string]interface{}, sort map[string]bool, index, _type string, sizeStr string) ([]byte, error) {
	var (
		list       []interface{}
		ess        *elastic.SearchService
		result     *elastic.SearchResult
		err        error
		object     interface{}
		byteResult []byte
		size       int
	)
	ess = slave.Search(index).Type(_type)
	if regexp != nil {
		for k, v := range regexp {
			ess = ess.Query(elastic.NewMatchQuery(k, v).Fuzziness("AUTO").Operator("AND"))
		}
	}
	if postFilter != nil {
		for k, v := range postFilter {
			ess = ess.PostFilter(elastic.NewTermQuery(k, v))
		}
	}
	if sort != nil {
		for k, v := range sort {
			if v {
				ess = ess.SortBy(elastic.NewFieldSort(k).Asc())
			} else {
				ess = ess.SortBy(elastic.NewFieldSort(k).Desc())
			}
		}
	}
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err == nil {
			ess = ess.Size(size)
		}
	}
	//ess = ess.Source([]string{"nickname", "usertype"})
	result, err = ess.Do(ctx)
	if err != nil {
		return nil, err
	}
	list = []interface{}{}
	if result != nil && result.Hits.TotalHits > 0 {
		for _, hit := range result.Hits.Hits {

			err := json.Unmarshal(*hit.Source, &object)
			if err != nil {
				return nil, err
			}
			list = append(list, object)
		}
	} else {
		return nil, nil
	}
	byteResult, err = json.Marshal(list)
	if err != nil {
		return nil, err
	}
	return byteResult, nil
}

// @Title Insert
// @Description update document by id
// @Parameters
//            object           map[string]interface{}       update map
//            index            string                       es index
//            _type            string                       es type
//            id               string                       es _id
// @Returns err:error
func (*es) Update(object map[string]interface{}, index, _type, id string) error {
	res, err := master.Update().
		Index(index).
		Type(_type).
		Id(id).
		DocAsUpsert(true).
		Doc(object).
		Refresh("true").
		Do(ctx)
	if err != nil {
		fmt.Println(res.Result)
		return err
	}
	return nil
}

// @Title Insert
// @Description insert document to es
// @Parameters
//            object           interface{}       insert struct object
//            index            string            es index
//            _type            string            es type
//            id               string            es _id
// @Returns err:error
func (*es) Insert(object interface{}, index, _type, id string) error {
	_, err := master.Index().
		Index(index).
		Type(_type).
		Id(id).
		BodyJson(object).
		Refresh("true").
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// @Title BulkInsert
// @Description insert documents from es by bulk
// @Parameters
//            objects          []interface{}     es bulk inert data
//            index            string            es index
//            _type            string            es type
// @Returns err:error
func (*es) BulkInsert(objects []interface{}, index, _type string) error {
	var (
		err error
		esb *elastic.BulkService
	)
	esb = master.Bulk()
	if objects != nil && len(objects) > 0 {
		for _, v := range objects {
			esb = esb.Add(elastic.NewBulkIndexRequest().Index(index).Type(_type).Id(strconv.Itoa(reflect.ValueOf(v).FieldByName("ID").Interface().(int))).Doc(v))
		}
		_, err = esb.Do(ctx)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf(ES_ERROR_DATA_EMPTY)
	}
	return nil
}

// @Title Delete
// @Description delete document from es by id
// @Parameters
//            index            string            es index
//            _type            string            es type
//            id               string            es _id
// @Returns err:error
func (*es) Delete(index, _type, id string) error {
	var (
		err error
		esd *elastic.DeleteService
	)
	esd = master.Delete().Index(index)
	if _type != "" {
		esd = esd.Type(_type)
	}
	if id != "" {
		esd = esd.Id(id)
	}
	_, err = esd.Refresh("true").Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

// @Title DeleteByQuery
// @Description delete document from es by query or id
// @Parameters
//            query            map[string]interface{}       query map
//            index            string                       es index
//            _type            string                       es type
// @Returns err:error
func (*es) DeleteByQuery(query map[string]interface{}, index, _type string) error {
	var (
		err   error
		esd   *elastic.DeleteByQueryService
		qlist []elastic.Query
		flag  = false
	)
	esd = master.DeleteByQuery(index).Type(_type)
	qlist = []elastic.Query{}
	if query != nil {
		for k, v := range query {
			if k == "all" {
				flag = true
				esd = esd.Query(elastic.NewMatchAllQuery())
				break
			} else {
				qlist = append(qlist, elastic.NewMatchQuery(k, v))
			}
		}
	}
	if !flag {
		esd = esd.Query(elastic.NewBoolQuery().Must(qlist...))
	}
	_, err = esd.Refresh("true").Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
