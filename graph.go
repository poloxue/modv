package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

var graphTemplate = `digraph {
node [shape=box];
{{ range $mod, $modId := .mods -}}
{{ $modId }} [label="{{ $mod }}"];
{{ end -}}
{{- range $modId, $depModIds := .dependencies -}}
{{- range $_, $depModId := $depModIds -}}
{{ $modId }} -> {{ $depModId }};
{{  end -}}
{{- end -}}
}
`

type ModuleGraph struct {
	Reader io.Reader

	Mods         map[string]int
	Dependencies map[int][]int
}

func NewModuleGraph(r io.Reader) *ModuleGraph {
	return &ModuleGraph{
		Reader: r,

		Mods:         make(map[string]int),
		Dependencies: make(map[int][]int),
	}
}

func (m *ModuleGraph) Parse() error {
	bufReader := bufio.NewReader(m.Reader)

	serialID := 1
	for {
		relationBytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		relation := bytes.Split(relationBytes, []byte(" "))
		mod, depMod := strings.TrimSpace(string(relation[0])), strings.TrimSpace(string(relation[1]))

		modId, ok := m.Mods[mod]
		if !ok {
			modId = serialID
			m.Mods[mod] = modId
			serialID += 1
		}

		depModId, ok := m.Mods[depMod]
		if !ok {
			depModId = serialID
			m.Mods[depMod] = depModId
			serialID += 1
		}

		if _, ok := m.Dependencies[modId]; ok {
			m.Dependencies[modId] = append(m.Dependencies[modId], depModId)
		} else {
			m.Dependencies[modId] = []int{depModId}
		}
	}
}

func (m *ModuleGraph) Render(w io.Writer) error {
	templ, err := template.New("graph").Parse(graphTemplate)
	if err != nil {
		return fmt.Errorf("templ.Parse: %v", err)
	}

	if err := templ.Execute(w, map[string]interface{}{
		"mods":         m.Mods,
		"dependencies": m.Dependencies,
	}); err != nil {
		return fmt.Errorf("templ.Execute: %v", err)
	}

	return nil
}
