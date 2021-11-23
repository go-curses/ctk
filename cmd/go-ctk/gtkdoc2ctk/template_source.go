package gtkdoc2ctk

var CtkSourceTemplate = `package <%= src.PackageName %>

import (
	"github.com/go-curses/cdk"
)

const Type<%= src.Name %> cdk.CTypeTag = "<%= src.PackageName %>-<%= src.Tag %>"

func init() {
	_ = cdk.TypesManager.AddType(Type<%= src.Name %>, func() interface{} { return <%= if (src.Constructor.String()) { %>Make<%= src.Name %>()<% } else { %>nil<% } %> })
}

<%= if (len(src.Hierarchy) > 0) { %><%= src.ObjectHierarchy() %><%= if (src.Description) { %>
<% } %><% } %><%= if (src.Description) { %><%= src.Description %><% } %>
type <%= src.Name %> interface {
	<%= src.Parent %><%= if (len(src.Implements) > 0) { %><%= for (impl) in src.Implements { %>
	<%= CamelCase(impl) %><% } %><% } %>

	Init() (already bool)
<%= for (f) in src.Functions { %>	<%= f.String() %>
<%  } %>}

// The C<%= src.Name %> structure implements the <%= src.Name %> interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with <%= src.Name %> objects.
type C<%= src.Name %> struct {
	C<%= src.Parent %>
}<%= if (src.Constructor.String()) { %>

// Make<%= src.Name %> is used by the Buildable system to construct a new <%= src.Name %>.
func Make<%= src.Name %>() *C<%= src.Name %> {
	return New<%= src.Name %>(<%= for (idx, arg) in src.Constructor.Argv { %><%= if (idx > 0) { %>, <% } %><%= arg.Value %><% } %>)
}

// <%= src.Constructor.Name %> is the constructor for new <%= src.Name %> instances.
func <%= src.Constructor.String() %> {
	<%= src.This %> := new(C<%= src.Name %>)
	<%= src.This %>.Init()
	return <%= src.This %>
}<% } %><%= if (len(src.Factories) > 0) { %><%= for (factory) in src.Factories { %>

<%= if (factory.Docs) { %><%= factory.Docs %>
<% } %>func <%= factory.String() %> {
	<%= src.This %> := new(C<%= src.Name %>)
	<%= src.This %>.Init()
	return <%= src.This %>
}<% } %><% } %>

// Init initializes an <%= src.Name %> object. This must be called at least once to
// set up the necessary defaults and allocate any memory structures. Calling
// this more than once is safe though unnecessary. Only the first call will
// result in any effect upon the <%= src.Name %> instance. Init is used in the
// New<%= src.Name %> constructor and only necessary when implementing a derivative
// <%= src.Name %> type.
func (<%= src.This %> *C<%= src.Name %>) Init() (already bool) {
	if <%= src.This %>.InitTypeItem(Type<%= src.Name %>, <%= src.This %>) {
		return true
	}
	<%= src.This %>.C<%= src.Parent %>.Init()
<%= if (len(src.Properties) > 0) { %><%= for (prop) in src.Properties { %>	_ = <%= src.This %>.InstallProperty(Property<%= prop.Name %>, <%= if (src.PackageName != "cdk") { %>cdk.<% } %><%= CamelCase(prop.Type.GoLabel) %>Property, <%= prop.Write %>, <%= sprintf("%v", prop.Default) %>)
<% } %><% } %>	return false
}
<%= if (len(src.Functions) > 0) { %><%= for (f) in src.Functions { %><%= if (f) { %><%= if (f.Docs) { %>
<%= f.Docs %><% } %>
func (<%= src.This %> *C<%= src.Name %>) <%= f.String() %> {<%= if (f.Body) { %>
<%= f.Body %>
<% } %>}
<% } %><% } %><% } %><%= if (len(src.Properties) > 0) { %><%= for (p) in src.Properties { %>
<%= p.Decl %>
<% } %><% } %><%= if (len(src.Signals) > 0) { %><%= for (s) in src.Signals { %>
<%= s.String() %><% } %><% } %>`
