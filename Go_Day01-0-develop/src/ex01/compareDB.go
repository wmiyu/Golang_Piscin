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

func nameInArr(newName string, oldNames []string) bool {
	for _, valueOld := range oldNames {
		if newName == valueOld {
			return true
		}
	}
	return false
}

func compareIngridients(ingredientsOld []XmlItem, ingredientsNew []XmlItem, cakeName string) {

	var flag bool = false

	for _, valueNew := range ingredientsNew {
		for _, valueOld := range ingredientsOld {
			if valueNew.Name == valueOld.Name {
				flag = true
				break
			}
		}
		if !flag {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", valueNew.Name, cakeName)
		}
		flag = false
	}

	for _, valueOld := range ingredientsOld {
		for _, valueNew := range ingredientsNew {
			if valueNew.Name == valueOld.Name {
				flag = true
				break
			}
		}
		if !flag {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", valueOld.Name, cakeName)
		}
		flag = false
	}

	for _, valueNew := range ingredientsNew {
		for _, valueOld := range ingredientsOld {
			if valueNew.Name == valueOld.Name {
				if valueOld.Unit == "" && valueNew.Unit != "" {
					fmt.Printf("ADDED unit \"%s\" for cake \"%s\"\n", valueNew.Unit, cakeName)
				}
				if valueOld.Unit != "" && valueNew.Unit != "" && valueOld.Unit != valueNew.Unit {
					fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", valueOld.Name, cakeName, valueNew.Unit, valueOld.Unit)
				}
				if valueOld.Count != valueNew.Count {
					fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", valueOld.Name, cakeName, valueNew.Count, valueOld.Count)
				}
				if valueOld.Unit != "" && valueNew.Unit == "" {
					fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", valueOld.Unit, valueOld.Name, cakeName)
				}
			}
		}
	}

}

func compare_dbs(resultOld []ResultRecipe, resultNew []ResultRecipe) {

	var oldNames []string
	var newNames []string

	for i := 0; i < len(resultOld); i++ {
		oldNames = append(oldNames, resultOld[i].Name)
	}

	for i := 0; i < len(resultNew); i++ {
		newNames = append(newNames, resultNew[i].Name)
	}

	for _, value := range newNames {
		if !nameInArr(value, oldNames) {
			fmt.Printf("ADDED cake \"%s\"\n", value)
		}
	}

	for _, value := range oldNames {
		if !nameInArr(value, newNames) {
			fmt.Printf("REMOVED cake \"%s\"\n", value)
		}
	}

	for _, nameOld := range resultOld {
		for _, nameNew := range resultNew {
			if nameOld.Name == nameNew.Name {
				if nameOld.Time != nameNew.Time {
					fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\" \n", nameOld.Name, nameNew.Time, nameOld.Time)
				}
				compareIngridients(nameOld.Ingridients, nameNew.Ingridients, nameOld.Name)
			}
		}
	}

}

func main() {

	flagOld := flag.String("old", "", "Old db file")
	flagNew := flag.String("new", "", "New db file")
	flag.Parse()

	var xmlFile, jsonFile DBReader = &Xml{}, &Json{}
	var resultOld, resultNew []ResultRecipe

	extensionOld := filepath.Ext(*flagOld)
	extensionNew := filepath.Ext(*flagNew)

	if extensionOld == ".xml" {
		resultOld = xmlFile.fileRead(*flagOld)
	} else if extensionOld == ".json" {
		resultOld = jsonFile.fileRead(*flagOld)
	}

	if extensionNew == ".xml" {
		resultNew = xmlFile.fileRead(*flagNew)
	} else if extensionNew == ".json" {
		resultNew = jsonFile.fileRead(*flagNew)
	}

	compare_dbs(resultOld, resultNew)
}
