package models

import (
	"fmt"
	"html/template"
	"reflect"
	"strconv"
)

// TODO Show should display decimals when required
// TODO figure out how to make Graphics a method of the implementation, not the interface

// Interface for all view types. Wrapped by the view struct to
// provide the 'Show' method, which compares a viewed field at
// the current stage of the simulation, with a compared field
// at a previous stage.
//
//	viewedField(f): returns the field f of a viewed record
//	comparedField(f): returns the field f of a compared record
type Viewer interface {
	ViewedField(f string) string
	ComparedField(f string) string
}

// Wrapper for the Viewer struct. Any view that implements Viewer
// can access the Show method of this type
type View struct {
	Viewer
}

// Returns a safe HTML string with a link to the ViewedField object
// Assumes the implementation supplies Name and ID fields
//
// v: an implementation of the Viewer interface
// urlBase: the root of the link url (eg `commodity`)
// template.HTML: safe string using ID and Name fields supplied by the implementation
func (v View) Link(urlBase string) template.HTML {
	return template.HTML(fmt.Sprintf("<a href=\"/%s/%s\">%s</a>", urlBase, v.ViewedField(`Id`), v.ViewedField(`Name`)))
}

// Returns a safe HTML string with a link to the Commodity of an industry
// Should be a method of  IndustryView but haven't yet figured out how
//
// v: an implementation of the Viewer interface
// template.HTML: safe string using ID and Name fields supplied by the implementation
func (v View) CommodityLink() template.HTML {
	//        <td><a href="/commodity/{{ .OutputCommodityId}}">{{ .Output }}</a> </td>
	return template.HTML(fmt.Sprintf(`<a href="/commodity/%s">%s</a>`, v.ViewedField(`OutputCommodityID`), v.ViewedField("Output")))
}

// Returns a safe HTML string with a graphic illustrating the origin
//
//	v: a CommodityView
//	template.HTML: safe string with a graphic representing the origin
func (v View) OriginGraphic() template.HTML {
	var htmlString template.HTML
	switch v.ViewedField(`Origin`) {
	case `INDUSTRIAL`:
		htmlString = "<i style=\"font-weight: bolder; color:blue\" class=\"fa fa-industry\"></i>"
	case `SOCIAL`:
		if v.ViewedField(`Usage`) == `Useless` {
			htmlString = "<i style=\"font-weight: bolder; color:rgba(128, 0, 128, 0.696)\" class=\"fas fa-user-tie\"></i>"
		} else {
			htmlString = "<i style=\"font-weight: bolder; color:red\" class=\"fa fa-user-friends\"></i>"
		}
	case `MONEY`:
		htmlString = "<i style=\"font-weight: 900; color:goldenrod\" class=\"fa fa-dollar\"></i>"
	default:
		htmlString = "Unknown Origin"
	}
	return template.HTML(htmlString)
}

// Returns a safe HTML string with a graphic illustrating the usage
//
//	v: a CommodityView
//	template.HTML: safe string with a graphic representing the usage
func (v View) UsageGraphic() template.HTML {
	var htmlString template.HTML
	switch v.ViewedField(`Usage`) {
	case `PRODUCTIVE`:
		htmlString = `<i style="font-weight: bolder; color:blue" class="fas fa-hammer"></i>`
	case `CONSUMPTION`:
		htmlString = `<i style="font-weight: bolder; color:green" class="fa fa-cutlery"></i>`
	case `MONEY`:
		htmlString = `<i class="fa fa-dollar" style="font-weight: 900; color:goldenrod"></i>`
	case `Useless`:
		htmlString = `<i class="fas fa-skull-crossbones" style="font-weight: bolder; color:black"></i>`
	default:
		htmlString = `Unknown Usage`
	}
	return template.HTML(htmlString)
}

// Provide a string, suitable for display in a template, that formats
// a viewed value and highlights values that have changed.
//
//	v: a View object
//	f: the name of the field to display
//	Returns: safe HTML string coloured red if the value has changed
func (v *View) Show(f string) template.HTML {
	vv, _ := strconv.Atoi(v.Viewer.ViewedField(f))
	vc, _ := strconv.Atoi(v.Viewer.ComparedField(f))
	var htmlString string
	if vv == vc {
		htmlString = fmt.Sprintf("<td style=\"text-align:center\">%d</td>", vv)
	} else {
		htmlString = fmt.Sprintf("<td style=\"text-align:center; color:red\">%d</td>", vv)
	}
	return template.HTML(htmlString)
}

// implements View for the Commodity Object
type CommodityView struct {
	viewedRecord   *Commodity
	comparedRecord *Commodity
}

// Provides the value of the field f in the viewedRecord of a CommodityView
//
//	f: the name of a field
//	c: a CommodityView
//	returns: the stringified value of the field (easiest generic solution)
func (c *CommodityView) ViewedField(f string) string {
	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(c.viewedRecord)).FieldByName(f))
}

// Provides the value of the field f in the comparedRecord of a CommodityView
//
//	f: the name of a field
//	c: a CommodityView
//	returns: the stringified value of the field (easiest generic solution)
func (c *CommodityView) ComparedField(f string) string {
	return fmt.Sprint(reflect.Indirect(reflect.ValueOf(c.comparedRecord)).FieldByName(f))
}

// Create a single CommodityView for display in a template
//
//	v: the currently viewed commodity
//	c: the same commodity at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateCommodityView(v *Commodity, c *Commodity) View {
	return View{&CommodityView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of CommodityView for display in a template
//
//	v: a slice of all commodities in the simulation at the current stage
//	c: a slice of the same commodities at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func CommodityViews(v *[]Commodity, c *[]Commodity) *[]View {
	var newViews = make([]View, len(*v))
	var vc *Commodity
	var cc *Commodity
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateCommodityView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

// Type for implementation of Viewer interface
type IndustryView struct {
	viewedRecord   *Industry
	comparedRecord *Industry
}

// Implements Viewer interface ViewedField method
func (i *IndustryView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *IndustryView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single IndustryView for display in a template
//
//	v: the currently viewed industry
//	c: the same industry at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryView(v *Industry, c *Industry) View {
	return View{&IndustryView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of IndustryView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryViews(v *[]Industry, c *[]Industry) *[]View {
	var newViews = make([]View, len(*v))
	var vc *Industry
	var cc *Industry
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateIndustryView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

// Experimental test to get at the stocks of an industry
func (ind IndustryView) GetVariable() View {
	fmt.Println("Entered GetVariable")
	viewedStock := ind.viewedRecord.Variable
	comparedStock := ind.comparedRecord.Variable
	return CreateIndustryStockView(viewedStock, comparedStock)
}

// Type for implementation of Viewer interface
type IndustryStockView struct {
	viewedRecord   *IndustryStock
	comparedRecord *IndustryStock
}

// Implements Viewer interface ViewedField method
func (i *IndustryStockView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ViewedField method
func (i *IndustryStockView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single IndustryStockView for display in a template
//
//	v: the currently viewed IndustryStock
//	c: the same IndustryStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateIndustryStockView(v *IndustryStock, c *IndustryStock) View {
	return View{&IndustryStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of IndustryStockViews for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func IndustryStockViews(v *[]IndustryStock, c *[]IndustryStock) *[]View {
	var newViews = make([]View, len(*v))
	var vc *IndustryStock
	var cc *IndustryStock
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateIndustryStockView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

// Type for implementation of Viewer interface
type NewClassView struct {
	viewedRecord   *Class
	comparedRecord *Class
}

// Implements Viewer interface ViewedField method
func (i *NewClassView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *NewClassView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single ClassView for display in a template
//
//	v: the currently viewed Class
//	c: the same Class at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassView(v *Class, c *Class) View {
	return View{&NewClassView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of ClassView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func NewClassViews(v *[]Class, c *[]Class) *[]View {
	var newViews = make([]View, len(*v))
	var vc *Class
	var cc *Class
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateClassView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}

// Type for implementation of Viewer interface
type NewClassStockView struct {
	viewedRecord   *ClassStock
	comparedRecord *ClassStock
}

// Implements Viewer interface ViewedField method
func (i *NewClassStockView) ViewedField(f string) string {
	s := reflect.Indirect(reflect.ValueOf(i.viewedRecord)).FieldByName(f)
	return fmt.Sprint(s)
}

// Implements Viewer interface ComparedField method
func (i *NewClassStockView) ComparedField(f string) string {
	return reflect.Indirect(reflect.ValueOf(i.comparedRecord)).FieldByName(f).String()
}

// Create a single ClassStockView for display in a template
//
//	v: the currently viewed ClassStock
//	c: the same ClassStock at an earlier point in the simulation
//	returns: a View object to supply to templates
func CreateClassStockView(v *ClassStock, c *ClassStock) View {
	return View{&NewClassStockView{
		viewedRecord:   v,
		comparedRecord: c,
	}}
}

// Create a slice of ClassStockView for display in a template
//
//	v: a slice of all industries in the simulation at the current stage
//	c: a slice of the same industries at an earlier point in the simulation
//	returns: a pointer to a slice of View objects to supply to templates
func NewClassStockViews(v *[]ClassStock, c *[]ClassStock) *[]View {
	var newViews = make([]View, len(*v))
	var vc *ClassStock
	var cc *ClassStock
	for i := range *v {
		vc = &(*v)[i]
		cc = &(*c)[i]
		newView := CreateClassStockView(vc, cc)
		newViews[i] = newView
	}
	return &newViews
}
