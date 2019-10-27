# pptxgrep
Search pattern from `.pptx` (Microsoft Power Point format) files like `grep` command.

## Usage
```
Usage:
  pptxgrep [options] pattern pptx1 [pptx2 ... pptxN]

Version:
  0.0.1

Options:
  -color
        colorize matched pattern
  -version
        print version
```

### Example 
```
$ pptxgrep 'Sample' sample.pptx
sample.pptx:1:Sample PowerPoint file.
sample.pptx:2:This is Sample.
```

result format is

```
target_file_path:slide_page_number:matched string
```


## Installation
```
$ go get github.com/noissefnoc/pptxgrep
```

## License
MIT

## Author
Kota Saito