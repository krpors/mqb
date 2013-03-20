package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	QUEUE      = "QU"  // queuename
	TARGCLIENT = "TC"  // 0 = jms, 1 = mq
	CCSID      = "CCS" // ccsid
)

// Parses a single line into components
func ParseSingleLine(line string) (jndiname, refaddr, name, value string, err error) {
	if strings.HasPrefix(line, "#") {
		err = fmt.Errorf("line is a comment")
		return
	}
	splitup := strings.Split(line, "/")

	// bail out if the splitup of the thing is bollox
	if len(splitup) != 4 {
		err = fmt.Errorf("split line must have length of 4")
		return
	}

	jndiname = splitup[0]
	refaddr = splitup[2]
	nvpair := strings.Split(splitup[3], "=")
	if len(nvpair) != 2 {
		err = fmt.Errorf("nvpair split must have length of 2")
		return
	}

	name = nvpair[0]
	value = nvpair[1]

	return
}

// Creates a new definition.
func NewDefinition() Definition {
	d := Definition{}
	d.Name = ""
	d.PropertyMap = make(map[string]Property, 3)
	return d
}

// Definition, which is basically a jndi name with its properties.
type Definition struct {
	Name        string
	PropertyMap map[string]Property // key is a RefAddr, which is just a number (0, 1, 2, ...)
}

func (d Definition) Queue() string {
	return d.get(QUEUE)
}

func (d Definition) CCSId() string {
	return d.get(CCSID)
}

func (d Definition) TargClient() string {
	s := d.get(TARGCLIENT)
	if s == "0" {
		return "JMS"
	}
	return "MQ"
}

func (d Definition) get(tp string) string {
	for _, prop := range d.PropertyMap {
		if prop.Type == tp {
			return prop.Content
		}
	}
	return ""
}

// Updates the internal map of Properties. Adds if not present, updates map if present. It's simple really:
// the key looked up is the RefAddr, which is just a number. If the RefAddr is found, the property is updated,
// depending on the name given in the function. If the name is "Content", the Content value is updated, etc.
// If the RefAddr key is not found, a new Property instance is created and THEN updated.
// 
// The content of the PropertyMap thus will be something like (pseudo):
//
//    0: {Content=7,    Type=VER, Encoding=String}
//    1: {Content=1208, Type=CCS, Encoding=String}
//    ... etc
func (d *Definition) UpdatePropertyMap(addr, name, value string) {
	theProp, found := d.PropertyMap[addr]
	if !found {
		theProp = Property{}
	}

	switch name {
	case "Content":
		theProp.Content = value
		break
	case "Type":
		theProp.Type = value
		break
	case "Encoding":
		theProp.Encoding = value
	}

	d.PropertyMap[addr] = theProp
}

// String representation of a Definition.
func (d Definition) String() string {
	return fmt.Sprintf("%s contains %v", d.Name, d.PropertyMap)
}

// Sortable definitions. Implements sort.Interface. Just sorts on name.
type Definitions []Definition

func (d Definitions) Len() int {
	return len(d)
}
func (d Definitions) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
func (d Definitions) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

// A property as serialized by the Websphere MQ tooling.
type Property struct {
	Content  string
	Encoding string
	Type     string
}

// Parses lines, returns a slice of Definitions.
func ParseLines(lines []string) Definitions {
	defmap := make(map[string]Definition)

	for _, line := range lines {
		jndiname, refaddr, name, value, err := ParseSingleLine(line)
		if err == nil {
			def, defFound := defmap[jndiname]
			if !defFound {
				def = NewDefinition()
				def.Name = jndiname
			}

			def.UpdatePropertyMap(refaddr, name, value)
			defmap[jndiname] = def
		}
	}

	definitions := make(Definitions, 0)

	for _, v := range defmap {
		definitions = append(definitions, v)
	}

	return definitions
}

// Reads the file f, and returns the lines as a string slice. Error is returned
// whenever something fails, like opening the file etc. etc.
func GetLinesFromFile(f string) ([]string, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		line, prefix, err := reader.ReadLine()
		if prefix {
			// TODO: we got a partial line only. didnt fit the buffer, shouldn't happen with these binding files
		}
		if err != nil {
			break
		}

		lines = append(lines, string(line))
	}

	return lines, nil
}

// Point d'entrance.
func main() {
	lines, err := GetLinesFromFile("C:/jndibinding/.bindings")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sep := "|"

	defs := ParseLines(lines)
	sort.Sort(defs)
	for _, def := range defs {
		fmt.Printf("%s%s%s%s%s%s%s\n", def.Name, sep, def.Queue(), sep, def.TargClient(), sep, def.CCSId())
	}
}
