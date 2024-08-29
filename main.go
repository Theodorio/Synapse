package main

import (
	"net/http"

	"os/exec"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

type Process struct {
	ImageName   string
	PID         string
	SessionName string
	SessionNum  string
	MemoryUsage string
}

func main() {
	r := gin.Default()

	// Serve HTML template at root URL
	r.LoadHTMLFiles("index.html", "processes.html")
	r.StaticFS("./static", http.Dir("static"))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "index"})

	})

	r.GET("/processes", func(c *gin.Context) {
		processes, err := getProcesses()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "processes.html", gin.H{"processes.html": processes})

		t, err := template.New("processes").Parse(`
			<table class="table">
				<thead>
					<tr>
						<th>Image Name</th>
						<th>PID</th>
						<th>Session Name</th>
						<th>Session</th>
						<th>Memory Usage</th>
					</tr>
				</thead>
				<tbody>
					{{ range . }}
					<tr>
						<td>{{ .ImageName }}</td>
						<td>{{ .PID }}</td>
						<td>{{ .SessionName }}</td>
						<td>{{ .SessionNum }}</td>
						<td>{{ .MemoryUsage }}</td>
					</tr>
					{{ end }}
				</tbody>
			</table>
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = t.Execute(c.Writer, processes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	})

	// Respond with JSON message at /ping endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server on port 1234
	err := r.Run(":1234")
	if err != nil {
		panic("[Error] Failed to start server due to:" + err.Error())
	}
}

func getProcesses() ([]Process, error) {
	cmd := exec.Command("tasklist", "/v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	processes := make([]Process, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		process := Process{
			ImageName:   fields[0],
			PID:         fields[1],
			SessionName: fields[2],
			SessionNum:  fields[3],
			MemoryUsage: fields[4],
		}

		processes = append(processes, process)
	}

	return processes, nil
}

func restartPC() {

}
