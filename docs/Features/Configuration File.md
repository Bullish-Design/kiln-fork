# Configuration file

The `kiln.yaml` file can be used instead of the flags to generate your site.

```
	Theme             string `yaml:"theme"`
	Font              string `yaml:"font"`
	URL               string `yaml:"url"`
	Name              string `yaml:"name"`
	Input             string `yaml:"input"`
	Output            string `yaml:"output"`
	Mode              string `yaml:"mode"`
	Layout            string `yaml:"layout"`
	FlatURLs          bool   `yaml:"flat-urls"`
	DisableTOC        bool   `yaml:"disable-toc"`
	DisableLocalGraph bool   `yaml:"disable-local-graph"`
	Port              string `yaml:"port"`
	Log               string `yaml:"log"`
```
