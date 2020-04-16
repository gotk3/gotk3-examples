package main

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"bufio"
	"regexp"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/glib"
)

type OptionsType struct {
	caseInsensitive bool
	wholeWord bool
	wholeLine bool
	filenameOnly bool
	filesWoMatches bool
}

type ResultsType struct {
	results	string
	err error
}

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Grep")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	if err != nil {
		log.Fatal("Unable to create mainBox:", err)
	}
	patternBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create patternBox:", err)
	}
	pathBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create pathBox:", err)
	}
	chkBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create chkBox:", err)
	}
	btnBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create btmBox:", err)
	}
	patternLbl, err := gtk.LabelNew("Pattern")
	if err != nil {
		log.Fatal("Unable to create patternLbl:", err)
	}
	patternEnt, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create patternEnt:", err)
	}
	pathLbl, err := gtk.LabelNew("Path")
	if err != nil {
		log.Fatal("Unable to create pathLbl:", err)
	}	
	pathEnt, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create pathEnt:", err)
	}
	caseInsensitiveChk, err := gtk.CheckButtonNewWithLabel("Case Insensitve")
	if err != nil {
		log.Fatal("Unable to create caseInsensitiveChk:", err)
	}	
	wholeWordChk, err := gtk.CheckButtonNewWithLabel("Whole Word")
	if err != nil {
		log.Fatal("Unable to create wholeWordChk:", err)
	}	
	wholeLineChk, err := gtk.CheckButtonNewWithLabel("Whole Line")
	if err != nil {
		log.Fatal("Unable to create wholeLineChk:", err)
	}	
	filenameOnlyChk, err := gtk.CheckButtonNewWithLabel("File Name Only")
	if err != nil {
		log.Fatal("Unable to create filenameOnlyChk:", err)
	}	
	filesWoMatchesChk, err := gtk.CheckButtonNewWithLabel("Files Without Matches")
	if err != nil {
		log.Fatal("Unable to create filesWoMatchesChk:", err)
	}	
	browseBtn, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Unable to create browseBtn:", err)
	}
	browseBtn.SetLabel("Browse")
	searchBtn, err := gtk.ButtonNew()
	if err != nil {
		log.Fatal("Unable to create searchBtn:", err)
	}
	statusLbl, err := gtk.LabelNew("")
	if err != nil {
		log.Fatal("Unable to create statusLbl:", err)
	}
	resultsTxt, err := gtk.TextViewNew()
	if err != nil {
		log.Fatal("Unable to create resultsTxt:", err)
	}
	resultsTxt.SetEditable(false)
	resultsTxt.SetRightMargin(80)
	resultsTxt.SetWrapMode(gtk.WRAP_WORD_CHAR)
	resultsBuf, err := resultsTxt.GetBuffer()
	if err != nil {
		log.Fatal("Unable to get buffer from resultsTxt:", err)
	}
	resultsSW, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal("Unable to create resultsSW:", err)
	}
	resultsSW.Add(resultsTxt)
	searchBtn.SetLabel("Search")
	patternBox.PackStart(patternLbl, false, false, 0)
	patternBox.PackStart(patternEnt, true, true, 0)
	pathBox.PackStart(pathLbl, false, false, 0)
	pathBox.PackStart(pathEnt, true, true, 0)
	chkBox.PackStart(caseInsensitiveChk, true, true, 0)
	chkBox.PackStart(wholeWordChk, true, true, 0)
	chkBox.PackStart(wholeLineChk, true, true, 0)
	chkBox.PackStart(filenameOnlyChk, true, true, 0)
	chkBox.PackStart(filesWoMatchesChk, true, true, 0)
	btnBox.PackStart(browseBtn, false, false, 0)
	btnBox.PackStart(searchBtn, false, false, 0)
	btnBox.PackStart(statusLbl, false, false, 0)
	mainBox.PackStart(patternBox, false, false, 0)
	mainBox.PackStart(pathBox, false, false, 0)
	mainBox.PackStart(chkBox, false, false, 0)
	mainBox.PackStart(btnBox, false, false, 0)
	mainBox.PackStart(resultsSW, true, true, 0)
	win.Add(mainBox)

	browseBtn.Connect("clicked", func() {
		fileChooserDlg, err := gtk.FileChooserNativeDialogNew("Open", win, gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", "_Cancel")
		if err != nil {
			log.Fatal("Unable to create fileChooserDlg:", err)
		}
		response := fileChooserDlg.NativeDialog.Run()
		if gtk.ResponseType(response) == gtk.RESPONSE_ACCEPT {
			fileChooser := fileChooserDlg
			filename := fileChooser.GetFilename()
			pathEnt.SetText(filename)
		} else {
			cancelDlg := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", "No file was selected")
			cancelDlg.Run()
			cancelDlg.Destroy()
		}
	})

	searchBtn.Connect("clicked", func() {
		pattern, _ := patternEnt.GetText()
		path, _ := pathEnt.GetText()
		options := OptionsType {}

		if caseInsensitiveChk.GetActive() {
			options.caseInsensitive = true
		}

		if wholeWordChk.GetActive() {
			options.wholeWord = true
		}
		if wholeLineChk.GetActive() {
			options.wholeLine = true
		}
		if filenameOnlyChk.GetActive() {
			options.filenameOnly = true
		}
		if filesWoMatchesChk.GetActive() {
			options.filesWoMatches = true
		}
		if pattern == "" {
			noPatternDlg := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", "No pattern was entered")
			noPatternDlg.Run()
			noPatternDlg.Destroy()
			return
		}
		if path == "" {
			noPathDlg := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", "No path was entered")
			noPathDlg.Run()
			noPathDlg.Destroy()
			return
		}
		
		if options.caseInsensitive == true {
			pattern = "(?i)" + pattern
		}
		if options.wholeWord == true {
			pattern = `\b` + pattern + `\b`
		}
		if options.wholeLine == true {
			pattern = "^" + pattern + "$"
		}

		go func() {

			glib.IdleAdd(func() {
				resultsBuf.SetText("")
				statusLbl.SetText("Working...")
				searchBtn.SetSensitive(false)
			})

			results, err := walkDir(pattern, path, options)			
			if err != nil {
				glib.IdleAdd(func() {
					statusLbl.SetText("")
				})
				errMsg := fmt.Sprintf("Error: %s", err)			
				glib.IdleAdd(func() {
					badPathDlg := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, 							gtk.BUTTONS_OK, "%s", errMsg)
					badPathDlg.Run()
					badPathDlg.Destroy()
				})
				glib.IdleAdd(func() {
					searchBtn.SetSensitive(true)
				})
				return
			}

			resultsSize := len(results)
			if resultsSize > 0 {
				if resultsSize <= 1000000 {
					glib.IdleAdd(func() {
						searchBtn.SetSensitive(true)
						resultsBuf.SetText(results)
						statusLbl.SetText("Success")
					})
					if err != nil {
						log.Fatal("Unable to execute IdleAdd()")
					}
				} else {
					if resultsSize <= 10000000 { 
						results := results[:resultsSize]
						msg := fmt.Sprintf("Showing only the first block (%d) of results", len(results))
						glib.IdleAdd(func() {
							searchBtn.SetSensitive(true)
							statusLbl.SetText(msg)
							resultsBuf.SetText(results)
						})
						if err != nil {
						log.Fatal("Unable to execute IdleAdd()")
						}
					} else {
						errMsg := fmt.Sprintf("Results size (%d) is too large", resultsSize)
						glib.IdleAdd(func() {
							searchBtn.SetSensitive(true)
							statusLbl.SetText(errMsg)
						})
						if err != nil {
							log.Fatal("Unable to execute IdleAdd()")
						}
					}
				}
			} else {
				_, err := glib.IdleAdd(func() {
					searchBtn.SetSensitive(true)
					statusLbl.SetText("No results were found")
				})
				if err != nil {
					log.Fatal("Unable to execute IdleAdd()")
				}
			}
		}()
	})

	win.SetDefaultSize(800, 600)
	win.ShowAll()
	gtk.Main()
}

func walkDir(pattern, dirToWalk string, options OptionsType) (string, error) {
	var matches []string
	err := filepath.Walk(dirToWalk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matchesFromFile, err2 := checkFileForPattern(path, pattern, options)
			if err2 != nil {
				log.Printf("Failed opening file: %s", err2)
			} else {
				matches = append(matches, matchesFromFile...)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result := strings.Join(matches, "")
	return result, nil
}

func isBinary(fileToRead string) (bool, error) {
	data := make([]byte, 256)
	file, err := os.Open(fileToRead)
	if err != nil {
		return false, err
	}
	defer file.Close()
	count, err := file.Read(data)
	if err != nil {
		return false, err
	}
	for i := 0; i < count; i++ {
		if data[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}

func checkFileForPattern(fileToRead string, pattern string, options OptionsType) ([]string, error) {
	matches := make([]string, 0)
  	r, err := regexp.Compile(pattern)
  	if err != nil {
		return nil, err
  	}
	file, err := os.Open(fileToRead)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return nil, nil
	}

	fileIsBinary, err := isBinary(fileToRead)
	if err != nil {
		return nil, err
	}	
	if fileIsBinary {
		log.Printf("%s is binary\n", fileToRead)
		return nil, nil	
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string 
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	if len(txtlines) == 0 {
		log.Printf("%s has no new line control characters.\n", fileToRead)
		return nil, nil
	}
	fileToRead = strings.ReplaceAll(fileToRead, `\\`, `\`)
	numNonMatches := 0
	for lineNum, line := range txtlines {
		if r.MatchString(line) {
			if options.filesWoMatches == false {
				if options.filenameOnly == true {
					match := fmt.Sprintf("%s\n\n", fileToRead)
					matches = append(matches, match)
					break				
				}
				var printableLine string
				var sb strings.Builder
				for _, r := range line {
					if int(r) >= 32 && int(r) != 127 {
						if r == '\\' || r == '"' {
							sb.WriteRune('\\')
						}
						sb.WriteRune(r)
					}
				}
				printableLine = sb.String()
				match := fmt.Sprintf("%s: %d:\n %s\n\n", fileToRead, lineNum + 1, printableLine)
				matches = append(matches, match)
			}
		} else {
			numNonMatches++
		}
	}
	if options.filesWoMatches == true && numNonMatches == len(txtlines) {
		match := fmt.Sprintf("%s\n\n", fileToRead)
		matches = append(matches, match)		
	}
	return matches, nil
}

