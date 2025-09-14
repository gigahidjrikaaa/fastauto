package match

import "testing"

func TestGlob(t *testing.T) {
    cases := []struct{ p, s string; ok bool }{
        {"main", "main", true},
        {"release/*", "release/v1", true},
        {"release/*", "hotfix/v1", false},
        {"feat?", "feat1", true},
        {"feat?", "feat12", false},
    }
    for _, c := range cases {
        if Glob(c.p, c.s) != c.ok {
            t.Fatalf("%s vs %s expected %v", c.p, c.s, c.ok)
        }
    }
}

