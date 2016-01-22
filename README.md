# Renderer

## Plans

So what a compiler is:

* Has 3-rd party storage interface
* Gets components from this storage
* Compiles components for rendering
* When asked for a component looks in storage
* Compiles components on demand
* Caches compiled results
* When component retrieved from storage is cached returns cached result

File system storage and watching for changes

* Components are very lightweight so full in-memory storage is fine along with file system replication
* Compiler itself does not have to watch for changes in files when a file system storage may be in-memory cached and on-change caches can be cleared.

Compiling & Rendering:

1. compiler -> storage [ask for component]
2. compiler [is component already compiled, in cache?]
3. compiler [if not cached: compile component]
4. compiler [if not cached: cache compiled component]
5. compiler [return compilation result]


### Components

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


### TODO

#### Components Storage

What we need is to create persistent storage interface with watching for changes.

Components persistent storage will be able to process and forward following event types:

	CreateFile // not component.json
	ChangeFile // not component.json
	RemoveFile // not component.json

	CreateComponent // component.json change
	ChangeComponent // component.json change
	RemoveComponent // component.json change


Compiler if enabled will recompile templates and store in-memory on event from persistent storage.
