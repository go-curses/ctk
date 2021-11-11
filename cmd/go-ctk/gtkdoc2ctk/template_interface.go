package gtkdoc2ctk

var CtkInterfaceTemplate = `package <%= src.PackageName %>

import (
	"github.com/go-curses/cdk"
)

<%= if (len(src.Hierarchy) > 0) { %><%= src.ObjectHierarchy() %><% } %><%= if (src.Description) { %>
<%= src.Description %><% } %>
type <%= src.Name %> interface {
<%= if (src.Parent == "CInterface") { %>	/* Base Interface */<% } else { %>	<%= src.Parent %>
<% } %>
<%= for (f) in src.Functions { %>
<%= f.InterfaceString() %>
<%  } %>}

<%= if (len(src.Properties) > 0) { %><%= for (p) in src.Properties { %>
<%= p.Decl %>
<% } %><% } %><%= if (len(src.Signals) > 0) { %><%= for (s) in src.Signals { %>
<%= s.String() %><% } %><% } %>`
