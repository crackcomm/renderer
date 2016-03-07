# Renderer

## Usage

### Installation

```sh
$ go get -u github.com/crackcomm/renderer
```

### API

Start API server:

```sh
$ renderer web -dir dashboard/
```

Render using API:

API responds with `JSON` **only when** `Accept` header contains `application/json`,
by default it writes HTML response.

```sh
# GET with component in URL
# http://127.0.0.1:6660/?name=dashboard.components&context={"components":[{"name":"example number one"}]}
# Or with POST method and component in body
$ curl -XPOST --data '{
  "name": "dashboard.components",
  "context": {
    "components": [
      {
        "name": "my.example.component"
      }
    ]
  }
}' http://127.0.0.1:6660/
```

### Components

Some example components can be found in `dashboard/components` directory.

JSON representation of an `example.root` component:

```json
{
  "name": "example.root",
  "main": "file://component.html",
  "styles": [
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css",
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css",
    "file://component.css"
  ],
  "scripts": [
    "https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js",
    "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js",
    "file://application.js"
  ]
}
```

Example of including components: `movies.movie` with `movies.video` component that can be rendered using `{{ video_component }}`:

Additionaly this component extends `movies.root` so after this component will be rendered, result will be included in `movies.root` under `{{ children }}` value.

```json
{
  "name": "movies.movie",
  "main": "file://component.html",
  "extends": "movies.root",
  "requires": {
    "video_component": {
      "name": "movies.video",
      "with": {
        "title": "{{ movie.title }} ({{ movie.year }})"
      }
    }
  }
}
```

### Rendering

Rendering `admin.domains` component with a list of `domains` in `context`.

```json
{
  "name": "admin.domains",
  "context": {
    "domains": [
      {"name": "test.pl", "links": ["test.pl"]},
      {"name": "test2.pl", "links": ["test.pl"]},
      {"name": "test3.pl", "links": ["test.pl"]}
    ]
  }
}
```

### Embedding

All data can be embed in `component.json`:

  * `http://` and `https://` urls are allowed
  * `file://` should be relative to `component.json`
  * `text://` for plain text (non templates)
  * `template://` will be executed as template

```json
{
  "name": "example.root",
  "main": "template://<h1>{{title}}</h1>",
  "styles": [
    "template://h1 { color: {{ color }}; }",
    "text://some text here"
  ],
  "scripts": [
    "template://console.log('{{ message }}');",
    "text://console.log('test');"
  ]
}
```


#### Some notes

So what a compiler is:

* Has 3-rd party storage interface
* Gets components from this storage
* Compiles components for rendering
* When asked for a component looks in storage
* Compiles components on demand
* Caches compiled results

File system storage and watching for changes

* Components are very lightweight so full in-memory storage is fine along with file system replication
* Compiler itself does not have to watch for changes in files when a file system storage may be in-memory cached and on-change caches can be cleared.

## License

Apache 2.0
