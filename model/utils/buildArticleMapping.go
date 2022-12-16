package utils

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"hellowiki/model/searchConfig"
)

func BuildArticleMapping(tokenOpt map[string]interface{}) mapping.IndexMapping {

	// a generic reusable mapping for english text
	standardIndexed := bleve.NewTextFieldMapping()
	standardIndexed.Store = false
	standardIndexed.IncludeInAll = false
	standardIndexed.Index = true
	standardIndexed.IncludeTermVectors = true
	standardIndexed.Analyzer = searchConfig.TokenName

	keywordIndexed := bleve.NewTextFieldMapping()
	keywordIndexed.Store = false
	keywordIndexed.IncludeInAll = false
	keywordIndexed.Index = true
	keywordIndexed.IncludeTermVectors = false
	keywordIndexed.Analyzer = "keyword"

	articleMapping := bleve.NewDocumentMapping()

	// title
	articleMapping.AddFieldMappingsAt("title", keywordIndexed)

	// content
	articleMapping.AddFieldMappingsAt("content", standardIndexed)

	// _all (disabled)
	disabledSection := bleve.NewDocumentDisabledMapping()
	articleMapping.AddSubDocumentMapping("_all", disabledSection)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("article", articleMapping)
	indexMapping.DefaultMapping = articleMapping

	indexMapping.AddCustomAnalyzer(searchConfig.TokenName, tokenOpt)
	indexMapping.DefaultAnalyzer = searchConfig.TokenName
	return indexMapping
}
