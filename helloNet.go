package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/kapacitor/pipeline"
	"github.com/influxdata/kapacitor/tick/stateful"
	"net/http"
	"time"
)

func main() {
	fmt.Println("server start...")
	http.HandleFunc("/", parseScriptHandler)
	http.ListenAndServe(":12345", nil)
}

func parseScriptHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	script := r.Form.Get("script")
	edgeType := r.Form.Get("edge-type")
	fmt.Println("script: ", script, "\n", "edge-type: ", edgeType)

	result := parseScript(script, edgeType)
	w.Write(result)
}

func parseScript(rawScript string, edgeType string) (cookedScript []byte) {
	d := deadman{}
	scope := stateful.NewScope()
	var p *pipeline.Pipeline
	var err error
	switch edgeType {
	case "stream":
		p, err = pipeline.CreatePipeline(rawScript, pipeline.StreamEdge, scope, d, nil)
	case "batch":
		p, err = pipeline.CreatePipeline(rawScript, pipeline.BatchEdge, scope, d, nil)
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("pipeline:", p)
	got, _ := json.MarshalIndent(p, "", "    ")
	return got
}

type deadman struct {
	interval  time.Duration
	threshold float64
	id        string
	message   string
	global    bool
}

func (d deadman) Interval() time.Duration { return d.interval }
func (d deadman) Threshold() float64      { return d.threshold }
func (d deadman) Id() string              { return d.id }
func (d deadman) Message() string         { return d.message }
func (d deadman) Global() bool            { return d.global }
