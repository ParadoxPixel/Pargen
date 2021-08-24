### Introducation
Pargen is a simple html page parser for go. By default it uses `%$1%` tags to identify parts to run through the parser and swapping it out for whatever the parsing provides.

### Example
Simple example of the default implementation, by default it doesn't re-parse files
```go
gen, err := NewDefaultGen(
    "./files/in", //Folder with files to parse
    "./files/parsed", //Folder to output parsed files
    func(s string) string {
        return s
    }, //Parsing function
)
if err != nil {
    panic(err)
}

//Prepares the input and output folder(gen.Prepare())
//Parses the files in the input folder(gen.ParseAll())
if err = gen.Initialize(); err != nil {
    panic(err)
}

//Load the parsed files to a *template.Template
templates, err := gen.Load(nil) //You can pass a template.FuncMap here
```
Normal way to get a `Gen` instance
```go
gen, err := NewGen(
	"./files/in",
	"./files/parsed",
	true, //If you want to re-parse existing files
	parser.DefaultRegex, //Regex to match tags
	func(s string) string {
		return s
	},
)
```