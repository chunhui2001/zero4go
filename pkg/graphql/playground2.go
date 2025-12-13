package graphql

import (
	"html/template"
	"net/http"
	"net/url"
)

var page = template.Must(template.New("graphiql").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8"/>
    <meta name="robots" content="noindex"/>
    <meta name="referrer" content="origin"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <title>{{ .Title }}</title>

    <link rel="stylesheet" href="{{.StaticServer}}/RichMedias/vendor/graphiql-0.13.2.css"/>
    <style>
        body {
            height: 100vh;
            margin: 0;
            overflow: hidden;
        }

        .variable-editor-title a {
            color: #3366cc;
        }

        #splash {
            color: #333;
            display: flex;
            flex-direction: column;
            font-family: system, -apple-system, "San Francisco", ".SFNSDisplay-Regular", "Segoe UI", Segoe, "Segoe WP", "Helvetica Neue", helvetica, "Lucida Grande", arial, sans-serif;
            height: 100vh;
            justify-content: center;
            text-align: center;
        }

        svg {
            fill: black;
        }
    </style>

    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/jquery-3.5.1.min.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/es6-promise.auto.min.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/fetch.min.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/react.min.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/react-dom.min.js"></script>

    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/graphiql-0.13.2.min.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/subscriptions-transport-ws-0.9.18.js"></script>
    <script type="text/javascript" src="{{.StaticServer}}/RichMedias/vendor/graphiql-subscriptions-fetcher-0.0.2.js"></script>
</head>
<body>
<div id="splash">
    Loading&hellip;
</div>
<script type="text/javascript">
    const graphqlServer = '{{- if .GraphqlServer }}{{ .GraphqlServer }}{{- end }}';
    const proto = location.protocol.startsWith("https") ? "https" : "http";
    const graphqlEndpoint = graphqlServer ? graphqlServer : proto + "://" + location.host + {{ .Endpoint }};

    function graphQLFetcher(graphQLParams) {

        const headers = new Headers();

        headers.append('Accept', 'application/json');
        headers.append('Content-Type', 'application/json');
        headers.append('Cache', 'no-cache');

        if (graphQLParams.headers) {

            for (var _key in graphQLParams.headers) {
                headers.append(_key, graphQLParams.headers[_key]);
            }
        }

        return fetch(graphqlEndpoint, {
            method: 'post',
            // mode: 'no-cors',
            mode: 'cors',
            headers: headers,
            body: JSON.stringify(graphQLParams.params ? graphQLParams.params : graphQLParams),
            credentials: 'include'
            // credentials: 'same-origin' // the fix
        }).then(function (response) {
            return response.text();
        }).then(function (responseBody) {
            try {
                return JSON.parse(responseBody);
            } catch (error) {
                return responseBody;
            }
        })
        .catch((reason) => {
            console.log(reason);
        });

    }

    var subscriptionsFetcher = window.GraphiQLSubscriptionsFetcher.graphQLFetcher(null, graphQLFetcher);

    ReactDOM.render(
        React.createElement(GraphiQL, {
            fetcher: subscriptionsFetcher,
            onEditHeaders: function(headers) {

            }
        }),
        document.body, () => {
            var copiedElement = $(".topBarWrap > .topBar > .title > span").clone();
            $(".topBarWrap > .topBar > .title > span").remove();
            $(".topBarWrap > .topBar > .title").append($("<a href='/'></a>").append(copiedElement));
        }
    );

</script>
</body>
</html>
`))

type GraphiqlConfig struct {
	Title                string
	StoragePrefix        string
	Endpoint             string
	FetcherHeaders       map[string]string
	UiHeaders            map[string]string
	EndpointIsAbsolute   bool
	SubscriptionEndpoint string
	JsUrl                template.URL
	JsSRI                string
	CssUrl               template.URL
	CssSRI               string
	ReactUrl             template.URL
	ReactSRI             string
	ReactDOMUrl          template.URL
	ReactDOMSRI          string
	EnablePluginExplorer bool
	PluginExplorerJsUrl  template.URL
	PluginExplorerJsSRI  string
	PluginExplorerCssUrl template.URL
	PluginExplorerCssSRI string

	GraphqlServer string
	StaticServer  string
}

// Playground Handler responsible for setting up the playground
func Playground(title, endpoint string) http.HandlerFunc {
	data := GraphiqlConfig{
		Title:              title,
		Endpoint:           endpoint,
		EndpointIsAbsolute: endpointHasScheme(endpoint),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")

		if err := page.Execute(w, data); err != nil {
			panic(err)
		}
	}
}

type GraphiqlConfigOption func(*GraphiqlConfig)

func WithGraphiqlFetcherHeaders(headers map[string]string) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.FetcherHeaders = headers
	}
}

func WithGraphiqlUiHeaders(headers map[string]string) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.UiHeaders = headers
	}
}

func WithGraphiqlVersion(jsUrl, cssUrl, jsSri, cssSri string) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.JsUrl = template.URL(jsUrl)
		config.CssUrl = template.URL(cssUrl)
		config.JsSRI = jsSri
		config.CssSRI = cssSri
	}
}

func WithGraphiqlReactVersion(
	reactJsUrl, reactDomJsUrl, reactJsSri, reactDomJsSri string,
) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.ReactUrl = template.URL(reactJsUrl)
		config.ReactDOMUrl = template.URL(reactDomJsUrl)
		config.ReactSRI = reactJsSri
		config.ReactDOMSRI = reactDomJsSri
	}
}

func WithGraphiqlPluginExplorerVersion(jsUrl, cssUrl, jsSri, cssSri string) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.PluginExplorerJsUrl = template.URL(jsUrl)
		config.PluginExplorerCssUrl = template.URL(cssUrl)
		config.PluginExplorerJsSRI = jsSri
		config.PluginExplorerCssSRI = cssSri
	}
}

func WithGraphiqlEnablePluginExplorer(enable bool) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.EnablePluginExplorer = enable
	}
}

func WithStoragePrefix(prefix string) GraphiqlConfigOption {
	return func(config *GraphiqlConfig) {
		config.StoragePrefix = prefix
	}
}

//// HandlerWithHeaders sets up the playground.
//// fetcherHeaders are used by the playground's fetcher instance and will not be visible in the UI.
//// uiHeaders are default headers that will show up in the UI headers editor.
//func HandlerWithHeaders(
//	title, endpoint string,
//	fetcherHeaders, uiHeaders map[string]string,
//) http.HandlerFunc {
//	return Playground(
//		title,
//		endpoint,
//		WithGraphiqlFetcherHeaders(fetcherHeaders),
//		WithGraphiqlUiHeaders(uiHeaders),
//	)
//}

// endpointHasScheme checks if the endpoint has a scheme.
func endpointHasScheme(endpoint string) bool {
	u, err := url.Parse(endpoint)
	return err == nil && u.Scheme != ""
}

// getSubscriptionEndpoint returns the subscription endpoint for the given
// endpoint if it is parsable as a URL, or an empty string.
func getSubscriptionEndpoint(endpoint string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		return ""
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	default:
		u.Scheme = "ws"
	}

	return u.String()
}
