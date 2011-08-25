// Copyright 2010 Abiola Ibrahim <abiola89@gmail.com>. All rights reserved.
// Use of this source code is governed by New BSD License
// http://www.opensource.org/licenses/bsd-license.php
// The content and logo is governed by Creative Commons Attribution 3.0
// The mascott is a property of Go governed by Creative Commons Attribution 3.0
// http://creativecommons.org/licenses/by/3.0/

package util

import (
	"container/vector"
	"os"
	"strings"
	"io/ioutil"
	"path"
	"fmt"
)
//StringBuilder type like StringBuilder in java
type StringBuilder struct {
	String string
}
//Creates an instance of StringBuilder
func NewStringBuilder(s string) *StringBuilder {
	sBuilder := new(StringBuilder)
	sBuilder.String = s
	return sBuilder
}
//Retrieves the string Content
func (this *StringBuilder) Content() string {
	return this.String
}
//Appends string to the end of the string content
func (this *StringBuilder) Append(s string) {
	this.String = this.String + s
}
//Deletes string from start index to end
func (this *StringBuilder) Delete(start, end int) {
	part1 := this.String[0:start]
	part2 := this.String[end:]
	this.String = part1 + part2
}
//Deletes the remaining string from index start
func (this *StringBuilder) DeleteTillEnd(start int) {
	this.String = this.String[0:start]
}
//Empties the string content
func (this *StringBuilder) Reset() {
	this.String = ""
}
//Returns the index of the first occurence of a particular string
func (this *StringBuilder) Index(s string) int {
	return strings.Index(this.Content(), s)
}
//Returns the length of the string content
func (this *StringBuilder) Len() int {
	return len(this.Content())
}

//Returns the substring from start to end index
func (this *StringBuilder) Sub(start, end int) string {
	return this.Content()[start:end]
}
//Returns the remaining string from the start index
func (this *StringBuilder) SubEnd(start int) string {
	return this.Content()[start:]
}
//QuoteParser to parse quotes e.g. { } or <?go ?>
type QuoteParser struct {
	buffer, static   *StringBuilder
	outer, inner     *vector.StringVector
	opening, closing string
}
//Creates a new QuoteParser with string s, opening and closing string
func NewQuoteParser(s, opening, closing string) *QuoteParser {
	parser := new(QuoteParser)
	parser.buffer, parser.static = NewStringBuilder(s), NewStringBuilder(s)
	parser.opening, parser.closing = opening, closing
	parser.inner, parser.outer = new(vector.StringVector), new(vector.StringVector)
	return parser
}
//Parses the string content in it
func (this *QuoteParser) Parse() (err os.Error) {
	for this.HasNext() {
		_, _, err = this.Next()
		if err != nil {
			return
		}
	}
	_, _, err = this.Next()
	return
}
//Returns the array of contents embedded in the quotes
func (this *QuoteParser) Parsed() []string {
	return []string(*this.inner)
}
//Returns the array of contents outside the quotes
func (this *QuoteParser) Outer() []string {
	return []string(*this.outer)
}
//Parses the next set and returns the embedded and outer strings with an error if any
//This method deletes the parsed string from the content
//If there is still need to parse whole content, use Reset()
func (this *QuoteParser) Next() (inner, outer string, err os.Error) {
	start := this.buffer.Index(this.opening)
	if start >= 0 {
		start += len(this.opening)
	}
	end := this.buffer.Index(this.closing)
	if end < 0 && start >= 0 {
		err = os.NewError("no closing string found near " + this.buffer.SubEnd(start-len(this.opening)))
		return
	}
	if this.HasNext() {
		inner = this.buffer.Sub(start, end)
	}
	if this.buffer.Len() > 0 {
		l := this.buffer.Index(this.opening)
		if l < 0 {
			l = this.buffer.Len()
		}
		outer = this.buffer.Sub(0, l)
		if start >= 0 {
			l = end + len(this.closing)
		}
		this.buffer.Delete(0, l)
		this.inner.Push(inner)
		this.outer.Push(outer)
	}
	return
}
//Resets the content to its state before parsing
func (this *QuoteParser) Reset() {
	this.buffer = this.static
}
//Checks whether there is next set of data to parse
func (this *QuoteParser) HasNext() (res bool) {
	start := this.buffer.Index(this.opening)
	end := this.buffer.Index(this.closing)
	return (start >= 0 && end >= 0)
}
//Returns the remaining content in the buffer being used
func (this *QuoteParser) String() string {
	return this.buffer.Content()
}
//public variable to store settings
var Config map[string][]string

const (
	SETTINGS = "pages.settings"
	ALL      = "ALL"
	NONE     = "NONE"
	EXCEPT   = "EXCEPT"
)
//Settings file type
type Settings struct {
	Data map[string][]string
}
//loads settings from pages.settings
func LoadSettings() (s *Settings, err os.Error) {
	s = new(Settings)
	err = s.parse()
fmt.Println("ss*****************************")
	return
}
//parse the informations in the settings file
func (this *Settings) parse() (err os.Error) {
	settings, err := ioutil.ReadFile(SETTINGS)
	if err != nil {
		return
	}
	parser := NewQuoteParser(string(settings), "{", "}")
	err = parser.Parse()
	if err != nil {
		return
	}
	key, values := parser.Outer(), parser.Parsed()
	this.Data = make(map[string][]string)
	for i, value := range values {
		this.Data[strings.TrimSpace(key[i])] = strings.Fields(value)
	}
	for _, s := range []string{"extensions", "handle", "srcfolder", "default"} {
		if this.Data[s] == nil {
			err = os.NewError("parsing settings file (pages.settings) failed : " + s + " invalid or not present")
			return
		}
	}
	fmt.Println("Before me")
	err = this.GeneratePages()
fmt.Println("my")

	if err != nil {
		return
	}
	vals := this.Data["handle"]
	if vals[0] == ALL {
		this.Data["handle"] = this.Data["pages"]
	} else if vals[0] == EXCEPT {
		vector := new(vector.StringVector)
		for _, page := range this.Data["pages"] {
			for i := 1; i < len(vals); i++ {
				if page != vals[i] {
					vector.Push(page)
				}
			}
		}
		this.Data["handle"] = []string(*vector)
	}
	return
}
//generates all .go source files
func (this *Settings) GeneratePages() (err os.Error) {
	if len(this.Data["extensions"]) == 0 {
		return
	}
	pages := new(vector.StringVector)
	fmt.Println("Before me7")
	err = this.iterFiles(this.Data["srcfolder"][0], pages)
	fmt.Println("Before me1")
	this.Data["pages"] = []string(*pages)
	return
}
//loops through root and subfolders to locate files with
//extensions specified in settings file
func (this *Settings) iterFiles(f string, pages *vector.StringVector) (err os.Error) {
	file, err := os.OpenFile(f, os.O_RDONLY, 0666)
	if err != nil {
		println(err.String())
		return
	}
	stat, er := file.Stat()
	if er != nil {
		err = er
		return
	}
	if stat.IsDirectory() {
		fmt.Println("iterFiles55555")
		dirs, err := file.Readdir(-1)
		if err != nil {
			return
		}
		for _, d := range dirs {
			this.iterFiles(path.Join(file.Name(), d.Name), pages)
		}
	} else {
		if hasExt(file.Name(), this.Data["extensions"]) {
			err = generate(file.Name())
			fmt.Println("iterFiles_eekkkk")
			if err != nil {
				return
			}
			pages.Push(file.Name())
		}
		file.Close()
	}
	return
}
//check if the file has the extension
func hasExt(filename string, ext []string) bool {
	extn := path.Ext(filename)
	for _, e := range ext {
		if "."+e == extn {
			return true
		}
	}
	return false
}

//direct function to generate .go source file
func generate(page string) (err os.Error) {
	p, err := NewPage(page)
	fmt.Println(p)
	fmt.Println(err)
	if err != nil {
		return
	}
	err = p.ParseToFile()
	fmt.Println("genterate_3")
	return
}
