# MOBI Reader

An implementation of a MOBI file reader written in Go.

## Usage

```
file, err := os.Open("data.mobi")
if err != nil {
	log.Fatal(err)
}
mobi, err := mobireader.Create(file)
if err != nil {
	log.Fatal(err)
}
```
