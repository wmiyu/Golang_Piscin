package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ResultRecipe struct {
	Name        string
	Time        string
	Ingridients []XmlItem
}

type Xml struct {
	Cake []XmlCake `xml:"cake"`
}

type XmlCake struct {
	Name        string         `xml:"name"`
	Time        string         `xml:"stovetime"`
	Ingridients XmlIngridients `xml:"ingredients"`
}

type XmlIngridients struct {
	Item []XmlItem `xml:"item"`
}

type XmlItem struct {
	Name  string `xml:"itemname" json:"ingredient_name"`
	Count string `xml:"itemcount" json:"ingredient_count"`
	Unit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

type Json struct {
	Cake []JsonCake `json:"cake"`
}

type JsonCake struct {
	Name        string    `json:"name"`
	Time        string    `json:"time"`
	Ingredients []XmlItem `json:"ingredients"`
}

type DBReader interface {
	fileRead(string) []ResultRecipe
}

func (x *Xml) fileRead(fileName string) []ResultRecipe {

	var xmlFile Xml
	result := []ResultRecipe{}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	err = xml.Unmarshal(byteValue, &xmlFile)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(xmlFile.Cake); i++ {
		result = append(result, ResultRecipe{Name: xmlFile.Cake[i].Name, Time: xmlFile.Cake[i].Time, Ingridients: xmlFile.Cake[i].Ingridients.Item})
	}

	return result
}

func (j *Json) fileRead(fileName string) []ResultRecipe {
	var jsonFile Json
	result := []ResultRecipe{}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	err = json.Unmarshal(byteValue, &jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(jsonFile.Cake); i++ {
		result = append(result, ResultRecipe{Name: jsonFile.Cake[i].Name, Time: jsonFile.Cake[i].Time, Ingridients: jsonFile.Cake[i].Ingredients})
	}

	return result
}

func main() {

	filePtr := flag.String("f", "", "a flag")
	flag.Parse()

	var xmlFile, jsonFile DBReader = &Xml{}, &Json{}
	var result []ResultRecipe

	extension := filepath.Ext(*filePtr)
	if extension == ".xml" {
		result = xmlFile.fileRead(*filePtr)
		data, _ := json.MarshalIndent(result, "", "    ")
		_ = ioutil.WriteFile("out_db.json", data, 0644)
	} else if extension == ".json" {
		result = jsonFile.fileRead(*filePtr)
		data, _ := xml.MarshalIndent(result, "", "    ")
		_ = ioutil.WriteFile("out_db.xml", data, 0644)
	}
}
