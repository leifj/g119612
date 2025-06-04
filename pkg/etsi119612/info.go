package etsi119612

import (
	"fmt"
	"slices"
	"errors"
	"github.com/subchen/go-xmldom"
)

func (tsl *TSL) schemeInfo() *xmldom.Node {
  	root := tsl.Doc.Root
	return root.QueryOne("SchemeInformation")
}

func (tsl *TSL) PrintInfo() {
	info := tsl.schemeInfo()
        v := info.FindOneByName("TSLVersionIdentifier")
        s := info.FindOneByName("TSLSequenceNumber")
        t := info.FindOneByName("TSLType")     
        fmt.Printf("%s (%s) type %s\n",v.Text,s.Text,t.Text)
}

func ServiceName(service *xmldom.Node,lang string) (string,error) {
	n := service.FindOneByName("ServiceName")
	names := n.FindByName("Name")
	idx := slices.IndexFunc(names,func(name *xmldom.Node) bool { return name.GetAttributeValue("lang") == lang });
	if (idx != -1) {
		return names[idx].Text,nil
	} else {
		return "",errors.New("no such language")
	}
}
