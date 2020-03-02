package main
//aaa
import (
	"bytes"
	"fmt"
	"github.com/influxdata/kapacitor/pipeline"
	"github.com/influxdata/kapacitor/pipeline/tick"
	"github.com/influxdata/kapacitor/tick/stateful"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("server start...")
	http.HandleFunc("/", parseScriptHandler)
	http.HandleFunc("/ok", parsePipelineHandler)
	http.ListenAndServe(":12345", nil)

}

// 脚本转为pipeline json
func parseScriptHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/x-www-form-urlencoded")
	w.WriteHeader(http.StatusCreated)

	edgeType := r.Form.Get("edge-type")
	script, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("script: ", script, "\n", "edge-type: ", edgeType)

	result := parseScript(string(script), edgeType)
	w.Write(result)
}

// pipeline json 转为脚本
func parsePipelineHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/x-www-form-urlencoded")
	w.WriteHeader(http.StatusCreated)

	pipelineJson := r.Form.Get("pipelineJson")

	// json to pipeline
	pipeline := pipeline.CreatePipelineSources()
	err := pipeline.Unmarshal([]byte(pipelineJson))
	if err != nil {
		panic(err)
	}
	// pipeline to script
	revertScript, _ := PipelineTick(pipeline)
	fmt.Println("revertScript: ", revertScript)

	w.Write([]byte(revertScript))
}

func parseScript(rawScript string, edgeType string) (cookedScript []byte) {
	d := deadman{}
	scope := stateful.NewScope()
	var p *pipeline.Pipeline
	switch edgeType {
	case "stream":
		p, _ = pipeline.CreatePipeline(rawScript, pipeline.StreamEdge, scope, d, nil)
	case "batch":
		p, _ = pipeline.CreatePipeline(rawScript, pipeline.BatchEdge, scope, d, nil)
	}
	if p == nil {
		return []byte("error pipeline is nil")
	}

	// 原始pipeline
	fmt.Println("RAE pipeline: ", p)

	// 转为json的原始pipeline
	marshalJSON, _ := p.MarshalJSON()
	fmt.Println("pipeline MarshalJSON: ", string(marshalJSON))

	// 原始pipeline恢复为脚本
	revertScript, _ := PipelineTick(p)
	fmt.Println("revert script: ", revertScript)

	// json to pipeline
	p2 := pipeline.CreatePipelineSources()
	err := p2.Unmarshal(marshalJSON)
	if err != nil {
		return []byte("error pipeline2 is nil")
	}
	// pipeline again to json
	marshalJSON2, _ := p2.MarshalJSON()
	fmt.Println("marshalJSON2 : ", string(marshalJSON2))

	return marshalJSON
}

func PipelineTick(pipe *pipeline.Pipeline) (string, error) {
	ast := tick.AST{}
	err := ast.Build(pipe)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	ast.Program.Format(&buf, "", false)
	return buf.String(), nil
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


func hello() (string) {
	return "hello"
}