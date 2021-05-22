# txtsearch

Search for text files.  
The search results can be listed or output to a directory.

## Usage

```
$ txtsearch cat dog -i .
total: 10  match: 3
pet.txt
cat.txt
dog.txt
```

The arguments are as follows.

```
Usage: txtsearch [flags] WORDS...
flags
  -i, --input string    Input directory. Specify the directory where the text files to be searched are located.
  -o, --output string   (optional) The directory to output text files matched by the search.
  -c, --cs              (optional) The search is case-sensitive. By default, it is not case-sensitive.
  -h, --help            Help.
```

Extracts files that contain any of the specified words.

You can output the files matched by `-o` to the specified directory.

```
$ txtsearch cat dog -i . -o out
total: 10  match: 3
copy the matched files to out

$ ls -1 out
pet.txt
cat.txt
dog.txt
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/txtsearch/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)

