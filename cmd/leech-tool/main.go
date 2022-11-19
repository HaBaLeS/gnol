package main

import (
	"encoding/json"
	"fmt"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/engine"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/modules"
	"github.com/HaBaLeS/gnol/cmd/leech-tool/persistence"
	"os"
	"time"
)

func main() {

	sessions := readSessionsJson("leech-data/iro.json")
	for _, s := range sessions {
		sf := s.LoadScrapeStatusFile()
		if sf != nil {
			if sf.LastScrapeTime.Add(time.Hour * 24).After(time.Now()) {
				fmt.Printf("Skipping %s , last scrate was within 24h\n", s.Name)
				continue
			}
		}
		s.Plm = &modules.Generic{
			NextSelector:  s.NextSelector,
			ImageSelector: s.ImageSelector,
			StopOnURl:     s.StopOnURl,
		}
		e := engine.Engine{
			Session: s,
		}
		e.Leech()
		s.WriteScrapeStatusFile()
		s.WriteMetaFile()
	}
}

func readSessionsJson(file string) []*persistence.Session {
	f, err := os.Open(file)
	panicIfErr(err)
	dec := json.NewDecoder(f)

	sessions := make([]*persistence.Session, 0)
	if err := dec.Decode(&sessions); err != nil {
		panic(err)
	}

	return sessions
}

func panicIfErr(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func configIRO() *persistence.Session {
	s := &persistence.Session{
		Count:   0,
		Start:   "https://irovedout.com/comic/iro/",
		Workdir: "iro2",
		Plm:     &modules.IROModule{},
	}
	return s
}

func configOglaf() *persistence.Session {
	s := &persistence.Session{
		Count:   0,
		Start:   "https://www.oglaf.com/cumsprite/",
		Workdir: "oglaf2",
		Plm:     &modules.OglafModule{},
	}
	return s
}

func configChester5000() *persistence.Session {
	//
	s := &persistence.Session{
		Count:   0,
		Start:   "http://jessfink.com/Chester5000XYV/?p=34",
		Workdir: "chester5000",
		Plm:     &modules.Chester5000Module{},
	}
	return s
}

func configCummoner() *persistence.Session {
	s := &persistence.Session{
		Count:   0,
		Start:   "http://totempole666.com/comic/first-time-for-everything-00-cover/",
		Workdir: "cummoner",
		Plm: &modules.Generic{
			NextSelector:  "a.comic-nav-next",
			ImageSelector: "div#comic a img",
		},
	}
	return s
}
