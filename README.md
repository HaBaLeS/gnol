# Gnol - a web based library for your Graphic Novel archives

## Features
- Browser based viewer for Desktop and Mobile
- Support wellknown formats cbz, cbr
- On-Demand loading of pages
- Only 1 executable, simple configuration
- Server-Side Image preprocessing to minimize traffic

## Planned Features
- PDF Support
- Multi-User Support
- Parental Control
- Save/Restore reading progress
- Search/Tag Comics/Novels

# Technology
- Go Backend
- Vanilla JavaScript
- Bulma CSS

## Running
install via go get

    go get github.com/HaBaLeS/gnol
    
run ./gnol

## Building from Source
checkout project, 

run make 

## Configuration
By default gnol will run on localhost:8666 it will scan for comic archives in working directory. You can use 

    ./gnol -c config.cfg 
    
to run gnol and provide a config file. Check example.cfg for details on configuration