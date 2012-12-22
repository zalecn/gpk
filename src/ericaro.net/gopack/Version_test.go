package gopack

import (
	"testing"
)


type Pair struct {
	s string
	v Version
}

func TestVersionParsing(t *testing.T) {
	
	pairs := []Pair{
		Pair{ "master"                , Version{0,0,0,"master",""		}},
		Pair{ "10.0.0"                , Version{10,0,0,"",""			   }},
		Pair{ "1.0.0"                 , Version{1,0,0,"",""				   }},
        Pair{ "1.0.0-alpha"           , Version{1,0,0,"alpha",""	       }},
        Pair{ "1.0.0-alpha.1         ", Version{1,0,0,"alpha.1",""         }},
        Pair{ "1.0.0-beta.2          ", Version{1,0,0,"beta.2",""          }},
        Pair{ "1.0.0-beta.11         ", Version{1,0,0,"beta.11",""         }},
        Pair{ "1.0.0-rc.1            ", Version{1,0,0,"rc.1",""            }},
        Pair{ "1.0.0-rc.1+build.1    ", Version{1,0,0,"rc.1","build.1"     }},
        Pair{ "1.0.0                 ", Version{1,0,0,"", ""               }},
        Pair{ "1.0.0+0.3.7           ", Version{1,0,0,"","0.3.7"           }},
        Pair{ "1.3.7+build           ", Version{1,3,7,"","build"           }},
        Pair{ "1.3.7+build.2.b8f12d7 ", Version{1,3,7,"","build.2.b8f12d7" }},
        Pair{ "1.3.7+build.11.e0f985a", Version{1,3,7,"","build.11.e0f985a"}},
 
	}
	
	for _,v := range pairs {
		ver, err := ParseVersion(v.s)
		if err != nil {
			t.Fatalf("Cannot parse%v \n", v)
		}
		if ver != 	v.v {
	
			t.Fatalf("Result mismatch (expected then result) \n%v\n%v \n", v.v, ver)
		}
	}
	
	//Parser("1.0.0-b.498.alpha+r.12", t)

}
