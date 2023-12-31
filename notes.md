# Introduce collections like a parent for different series
- [ ] default collection is all without the NSFW series
- 

# Series
- [ ] Series Sorting
- [ ] series nsfw flag
- [ ] series tags
- [ ] Series Editor
- [ ] Change Cover

# Comic
- [ ] Change Cover
- [ ] Change first image

# General

- [ ] additional metadata with comic.xml
- [ ] FIXME Auth without encrypted cookie! Use SessionID and server side Session in DB

# Workflows

- Buy HumbleBundle Comic Bundle
- Download all PDF's
- Convert PDF into cbz (it does also some resizing and tagging)
- Create a new Series in the Gnol UI
- Upload the cbz to Gnol

## Humblebundle Download (pdf's only whole package)
Install [Humble Bundle Downloader](https://github.com/xtream1101/humblebundle-downloader)

    hbd  -i pdf -s "eyJyZXBsX3thars$23erwhst3sdShesdasdaDwasaW1lIjoxNzAwMzE2NTA0fQ==|1700316512|f4edfad4afb7be0d0d44e17e7c5cc15d683d42bd" --library-path "." -k A3umahAEMbp5C2t8

## Convert all PDF in all folders to cbz (With Tags)

    find . -name '*.pdf' -exec gnol-tools pdf2cbz --tags=Humblebundle,TopCow  -v {} \;

##  Upload all cbz files in all subfolders (into a series with id 36)

    find . -name '*.cbz' -exec gnol-tools upload -v --seriesId=36 {} \;

##  