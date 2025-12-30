package graphql

import (
	"fmt"
	"net/http"
)

// Handler serves the GraphiQL UI
func Handler(w http.ResponseWriter, r *http.Request) {
	theme := DetectTheme(r)
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>casspeed GraphQL API</title>
	<style>
		body {
			height: 100%%;
			margin: 0;
			width: 100%%;
			overflow: hidden;
		}
		#graphiql {
			height: 100vh;
		}
		%s
	</style>
	<script crossorigin src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
	<script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
	<link rel="stylesheet" href="https://unpkg.com/graphiql/graphiql.min.css" />
</head>
<body>
	<div id="graphiql">Loading...</div>
	<script src="https://unpkg.com/graphiql/graphiql.min.js" type="application/javascript"></script>
	<script>
		const fetcher = GraphiQL.createFetcher({
			url: '/graphql/query',
		});

		const root = ReactDOM.createRoot(document.getElementById('graphiql'));
		root.render(
			React.createElement(GraphiQL, {
				fetcher,
				defaultQuery: '# Welcome to casspeed GraphQL API\n\nquery {\n  health\n}\n',
			}),
		);
	</script>
</body>
</html>
`, GetThemeCSS(theme))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// QueryHandler handles GraphQL queries
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	// Simple placeholder GraphQL response
	// In a real implementation, use a GraphQL library like graphql-go
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	response := `{
	"data": {
		"health": "ok"
	}
}
`
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, response)
}
