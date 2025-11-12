package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	port := "3000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	// Serve static files
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// API endpoint to open profile
	http.HandleFunc("/api/open-profile", func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")
		profileType := r.URL.Query().Get("type")
		view := r.URL.Query().Get("view")

		if service == "" || profileType == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		profileFile := filepath.Join("profiles", fmt.Sprintf("%s-%s.prof", service, profileType))

		// Check if profile exists
		if _, err := os.Stat(profileFile); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("Profile not found: %s", profileFile), http.StatusNotFound)
			return
		}

		// Build pprof command
		args := []string{"tool", "pprof", "-http=:8080"}
		if view != "" {
			args = append(args, fmt.Sprintf("-%s", view))
		}
		args = append(args, profileFile)

		// Execute pprof in background
		cmd := exec.Command("go", args...)
		if err := cmd.Start(); err != nil {
			http.Error(w, fmt.Sprintf("Failed to start pprof: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"success","message":"Profile opened at http://localhost:8080"}`)
	})

	// API endpoint to list available profiles
	http.HandleFunc("/api/profiles", func(w http.ResponseWriter, r *http.Request) {
		profiles := []string{}

		if _, err := os.Stat("profiles"); !os.IsNotExist(err) {
			files, err := filepath.Glob("profiles/*.prof")
			if err == nil {
				for _, file := range files {
					profiles = append(profiles, filepath.Base(file))
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"profiles":%q}`, profiles)
	})

	// API endpoint to refresh profiles
	http.HandleFunc("/api/refresh-profiles", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Refreshing profiles...")

		// Determine which script to run
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", "profile.bat")
		} else {
			cmd = exec.Command("bash", "profile.sh")
		}

		// Run the profiling script
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to refresh profiles: %v\n%s", err, output)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"status":"error","message":"Failed to run profile script: %v"}`, err)
			return
		}

		log.Println("Profiles refreshed successfully")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"success","message":"Profiles refreshed successfully"}`)
	})

	url := fmt.Sprintf("http://localhost:%s/profile-viewer.html", port)
	fmt.Printf("🚀 Profile Viewer Server Starting...\n")
	fmt.Printf("📊 Dashboard: %s\n", url)
	fmt.Printf("📁 Serving from: %s\n\n", getCurrentDir())
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// Try to open browser
	go openBrowser(url)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		fmt.Printf("Please open: %s\n", url)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Please open: %s\n", url)
	}
}
