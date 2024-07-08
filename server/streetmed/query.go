package streetmed

import (
	"strings"
)

type query struct {
    q []string
}

func parseQuery(qs string) (query, error) {
    var q []string
    var sb strings.Builder

    quoted := false
    escaped := false
    for _, c := range strings.ToLower(qs) {
        switch c {
        case '\\':
            if !escaped {
                escaped = true
                continue
            }
        case '\'', '"':
            if !escaped {
                quoted = !quoted
                continue
            }
        case ' ':
            if !quoted {
                if toAdd := sb.String(); toAdd != "" {
                    q = append(q, sb.String())
                }
                sb.Reset()
                continue
            }         
        }

        if escaped { escaped = false }
        // only write lowercase ascii
        if int(c) > 96 && int(c) < 123 {
            sb.WriteRune(c)
        }
    }

    if toAdd := sb.String(); toAdd != "" {
        q = append(q, sb.String())
    }

    return query{ q }, nil
}

func (this query) LinearSearch() ([]set[string], error) {
    var results []set[string]
    for _, q := range this.q {
        tmp := make(set[string])
        for k, v := range fileMap {
            if strings.Contains(v, q) { tmp.Add(k) } 
        }
        results =  append(results, tmp)
    }

    return results, nil
}

func (this query) TreeSearch() ([]set[string], error) {
    var results []set[string]
    for _, q := range this.q {
        if q == "" { continue }
        tmp := treeConst.search(q)
        results = append(results, tmp)
    }

    return results, nil
}

func (this query) CmdlineString() string {
    var sb strings.Builder
    for i, s := range this.q  {
        if i != 0 {
            sb.WriteString(" and ")
        }
        sb.WriteString("'"+s+"'")
    }

    return sb.String()
}
