// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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