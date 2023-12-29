<!DOCTYPE html>
<html>
  <head>
    <title>Helm Charts</title>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.min.css" />
    <style>
      .markdown-body {
        box-sizing: border-box;
        min-width: 200px;
        max-width: 980px;
        margin: 0 auto;
        padding: 45px;
      }
      @media (max-width: 767px) {
        .markdown-body {
          padding: 15px;
        }
      }
    </style>
  </head>

  <body>

    <section class="markdown-body">
      <h1>Helm Charts</h1>

      <h2>Usage</h2>
      <pre lang="no-highlight"><code>
        helm repo add haproxy-operator https://six-group.github.io/haproxy-operator
      </code></pre>

      <h2>Charts</h2>

      <ul>
			{{range $key, $chartEntry := .Entries }}
				<li>
					<p>
						{{ (index $chartEntry 0).Name }}
						(<a href="{{ (index (index $chartEntry 0).Urls 0) }}" title="{{ (index (index $chartEntry 0).Urls 0) }}">
						{{ (index $chartEntry 0).Version }}
						@
						{{ (index $chartEntry 0).AppVersion }}
						</a>)
					</p>
					<p>
						{{ (index $chartEntry 0).Description }}
					</p>
				</li>
			{{end}}
      </ul>
      <p>
		<time datetime="{{ .Generated.Format "2006-01-02T15:04:05" }}" pubdate id="generated">{{ .Generated.Format "Mon Jan 2 2006 03:04:05PM MST-07:00" }}</time>
      </p>
    </section>
  </body>
</html>