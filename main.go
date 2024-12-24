/********************************************************************************
* La funzione iterara su ogni file e carpella al interno di PATH e fare due cose:
* 1. limpiara i contenuti dei file di tutti i id uniche forniti da notion
* 2. Rimovera la id unica dal nome dil file.
* ******************************************/

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// DIR
const PATH = "./notion_12.22.24"
const ITER_LOOPS = 10

func main(){ 

	// So che questo potrebbe essere megliorato. Non c'e bisogna di fare una loop
	// su un valore fisso quando potrei simplicemente calcolare la quantita di
	// loops essata derivata da getAllDirectoriesToRename(). Purtroppo, sono
	// stanco e ho altri proggeti di fare. Fore lo faro nel futuro. 
	// for index, _ := range [ITER_LOOPS]int{1:10} {
	// 	if index < ITER_LOOPS {
	// 		err := cleanPaths(PATH)
	// 		fmt.Println(err)
	// 	}
	// }

	err := cleanNotionLinks(PATH)

	if err != nil {
		fmt.Println("********", err)
	}

	// bkp.Bkp()
}

// limpia i nomi di tutti i file trovati al'interno di path
func cleanPaths(path string) error {

	paths, err := getAllDirectoriesToRename(path)

	// restituici se non puoi riaggiungere i file dirs
	if err != nil{
		return err
	}

	for _, path := range paths {
		err = cleanSinglePath(path)
		if err != nil{
			fmt.Println(err)
		}
	}

	return err
}

// Will analyze the given path and will call the directory name cleaning function 
// or the file cleaning function
func cleanSinglePath(path string) error {
	// check that the path exists before cleaning
	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("🆘🔍 %s does not exists: %v \n", path, err)
		} else{
			return fmt.Errorf("❌📥 I could not access %s \n error: %v", path, err)
		}
	}

	// decide which cleaning function to use according to the path type
	if info.IsDir(){ 
		err = cleanDirName(path)
	} else {
		err = cleanFileName(path)
	}

	if err != nil {
		return err
	}

	return nil
}

//  will remove a notion unique id from dir name 
func cleanDirName(path string) error {

	//get the last part of the name, which is VERY likely an id
	splitPath := strings.Split(path, " ")
	var uniqueId  = ""

	if len(splitPath) > 1 {
		uniqueId = splitPath[len(splitPath) - 1]
	}

	// Contains a unique Id, remove it
	newPathWOId := strings.Replace(path, uniqueId, "", -1)
	// now remove any trailing " ", -, or _
	newPathWOId = strings.Trim(newPathWOId, " ")
	newPathWOId = strings.Trim(newPathWOId, "-")
	newPathWOId = strings.Trim(newPathWOId, "_")

	fmt.Printf("💈 Changing old name from %s -> %s \n", path, newPathWOId)

	err := os.Rename(path, newPathWOId)

	if err != nil {
		return fmt.Errorf("❌ %s I could not rename this directory: \n error: %v \n", path, err)
	}

	// everything went great 
	return nil
}

// will remove a notion unique id from file name 
func cleanFileName(path string ) error {
	// what is the extension of this file
	fileExt := filepath.Ext(path)

	//get the last part of the name, which will also contain the extension 
	splitPath := strings.Split(path, " ")
	var uniqueId  = ""

	if len(splitPath) > 1 {
		uniqueId = splitPath[len(splitPath) - 1]
	}

	// Contains a unique Id, remove it
	newNameWOId := strings.Replace(path, uniqueId, "", -1)

	// now remove any trailing " ", -, or _
	newNameWOId = strings.Trim(newNameWOId, " ")
	newNameWOId = strings.Trim(newNameWOId, "-")
	newNameWOId = strings.Trim(newNameWOId, "_")

	// rimette la estenzione
	newNameWOId+= fileExt
	

	// Does this file we are trying to rename, already exist?
	if isFileExists(newNameWOId) {
		return fmt.Errorf("File already exists. Skipping: %s", path)
	}

	fmt.Printf("💈 Changing old name from %s -> %s \n", path, newNameWOId)

	// all tests passed. Rename now
	err := os.Rename(path, newNameWOId)

	if err != nil {
		return fmt.Errorf("❌ There was an error renaming this file: %v", err)
	}

	return nil
}

// Verifica se il path dato essiste
func isFileExists(path string) bool{
	_, err :=os.Stat(path)
	return os.IsExist(err)
}

// Caminero su tutti i dile al'interno de una carpella path
func getAllDirectoriesToRename(path string) ([]string, error) {
	// empty array to append paths to
	var paths = make([]string, 0)

	// la posizione indice dil file 
	var index = 0

	// walk through the path to get all the directories to append 
	filepath.Walk(path, func (path string, info fs.FileInfo, err error) error {
		// ignora la carpella radice
		if index == 0 {
			index++
			return nil
		}

	

		// assicurarti che possa accessare il file 
		if err != nil {
			fmt.Printf("Errore nell'accesso al file %s: %v\n", path, err)
			return nil
		}

		// assicurarti che info non sia nil
		if info == nil {
			fmt.Printf("Info è nil per il percorso: %s\n", path)
			return nil
		}

		// Split the file name by words
		splitFileName := strings.Split(info.Name(), " ")

		// is this a file with more than one word? if yest get the last part
		var lastFileNamePart string
		if (len(splitFileName) > 1){
			lastFileNamePart = splitFileName[len(splitFileName)-1]
		} else {
			// if not, this file does not have an ID, you can come on
			fmt.Printf("✅ %s is not a dirty name \n", info.Name())
			return nil
		}

		// Prepare a regex that will check for a unique ID of more than 12 chars 
		// and with at least 3 numbers in it
		rgx := regexp.MustCompile(`(?:\D*\d){4}`)
		// is the last part of the file name a unique ID? 
		isMatch := rgx.MatchString(lastFileNamePart)

		// If not, log it but do not add it to the slice
		if isMatch {
			paths = append(paths, path)
		} else {
			fmt.Printf("✅ %s is not a dirty name \n", info.Name())
		}

		return nil
	})

	// log the dirty names at the end so I can more easily read them
	for _, path := range paths {
		fmt.Printf("❌ %s needs cleaning", path)
	}

	fmt.Printf("🔢 Total paths to clean %d", len(paths))

	return paths, nil
}


func getAllHTMLFilesInPath(path string) []string{

	var htmlFiles = make([]string, 0)

	filepath.Walk(path, func (path string, info os.FileInfo, err error) error {
		if filepath.Ext(info.Name()) == ".html"{
			htmlFiles = append(htmlFiles, path)
		}

		return nil
	})


	return htmlFiles
}

func cleanNotionLinks(path string) error {
    // Regular expression to match Notion's unique ID pattern
    // This matches strings that look like " 1234abcd" at the end of text
    	idDirPattern := regexp.MustCompile(`\b[a-zA-Z0-9]*\.?[a-zA-Z0-9]{16,}\b`)
   	idFilePattern := regexp.MustCompile(`^[a-zA-Z0-9]{16,}$`)
	totalLinksFound := 0

    // Walk through all files in the workspace
	for _, path := range getAllHTMLFilesInPath(path) {

		fmt.Printf(
		"=================================================\n" + 
		"====Starting file: %s==== \n" + 
		"================================================= \n",

		path,
	)

        // Skip if not an HTML file
        if !strings.HasSuffix(strings.ToLower(path), ".html") {
            return nil
        }

        // Open the file
        file, err := os.Open(path)
        if err != nil {
            return fmt.Errorf("error opening file %s: %v", path, err)
        }
        defer file.Close()

        // Parse HTML
        doc, err := html.Parse(file)
        if err != nil {
            return fmt.Errorf("error parsing HTML in %s: %v", path, err)
        }

        // Recursively process nodes and return  true if there were any changes that happened in
	  // this file so we can rewrite it
        modified := processHtmlNode(doc, idDirPattern, idFilePattern)
	
	  if !modified{
  		fmt.Println("❌", path, modified)
	  }else {
  		fmt.Println("✅", path, modified)
	  }

        // If modifications were made, write the changes back to the file
        if modified {
            // Create a temporary file
            tempFile, err := os.CreateTemp(filepath.Dir(path), "temp-*.html")
            if err != nil {
                return fmt.Errorf("error creating temp file for %s: %v", path, err)
            }
            tempPath := tempFile.Name()

		fmt.Printf("🪄🦄 Temp file %s created", tempPath)

            // Write the modified HTML to the temp file
            err = html.Render(tempFile, doc)
            tempFile.Close()
            if err != nil {
                os.Remove(tempPath)
                return fmt.Errorf("error writing to temp file for %s: %v", path, err)
            }

            // Replace the original file with the temp file
            err = os.Rename(tempPath, path)
            if err != nil {
                os.Remove(tempPath)
                return fmt.Errorf("error replacing file %s: %v", path, err)
            }

            fmt.Printf("✅✨ File Processed: %s\n", path)
        }
    }

//     if err != nil { 
// 	return fmt.Errorf("****** ERR: %v", err)
//     }

	// how many links were found?
	fmt.Printf("🔢 total links found: %d \n", totalLinksFound)
	return nil
}

// Receives parsed html page and will replace the id in all html <a> elements found matching according to
// idDirPattern and idFilePattern
func processHtmlNode(n *html.Node, idDirPattern *regexp.Regexp, idFilePattern *regexp.Regexp) bool{

	modified := false

	if n.Type == html.ElementNode {
		// Process href attributes in anchor tags
		if n.Data == "a" {
			for i, attr := range n.Attr {
			// Original href value 
			ogName := n.Attr[i].Val

			if attr.Key == "href" {

				// if the href is a link to a file and not a dir will have .html extension
				isHTMLFile := false

				// new name formed without the id 
				var newNameParts = make([]string, 0)

				// Check for an ID at the herf level by dir "/"
				// fmt.Printf("◉ All parts in href: %s \n", strings.Join(strings.Split(attr.Val, "/"), ", "))
				
				for _, part := range strings.Split(attr.Val, "/") {
					// fmt.Printf("  ◎ Processing part: %s \n", part)
					if idDirPattern.MatchString(part){
						// If the directory has and ID, then check by word level " "
						// fmt.Printf("    ❌ %s part is dirty! Checking by space now \n", part)
						
						nameParts := make([]string, 0)
						
						// fmt.Printf("      ◦ Parts by space: %s \n", strings.Join(strings.Split(part, " "), ", "))

						partsSplitBySpace := make([]string, 0)

						for _, part := range strings.Split(part, " "){
							splitPart := strings.Split(part, "%20")
							partsSplitBySpace = append(partsSplitBySpace, splitPart...)
						}

						for _, word := range partsSplitBySpace {
							// fmt.Printf("        ・Processing word: %s \n", word)

							// remove html from the string if this is a file file
							wordWithoutExtension, hasExt := strings.CutSuffix(word, ".html")
							isHTMLFile = hasExt

							if idFilePattern.MatchString(wordWithoutExtension) {
								//totalLinksFound++
								fmt.Printf("          ❌ %s is dirty! removing... \n", wordWithoutExtension)

								// Let the function know that the file has been updated
								modified = true
								
							} else{
								fmt.Printf("          ✅ %s is clean! skipping... \n", wordWithoutExtension)
								nameParts = append(nameParts, wordWithoutExtension)
							}
						}

						// Since this string was originally split by space, join them by space
						newName := strings.Join(nameParts, " ")

						// If it is a link to an html file, add the extension back
						if isHTMLFile {
							newName+= ".html"
						}

						// add this name to the newNameParts 
						newNameParts = append(newNameParts, newName)
					} else {
						// fmt.Printf("    ✅ %s is clean! skipping... \n", part)
						newNameParts = append(newNameParts, part)
					}
				}

				// since this was originally split by /, join them by /
				ogName = strings.Join(newNameParts, "/")
				n.Attr[i].Val = ogName

				// fmt.Printf("📝 New name will be:  %s \n", ogName)
			}

			// Now change the displayable text to match the href. 
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode{
					c.Data = ogName
				}
			}
			}
		}
	}

	// Get into each node of the html
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if processHtmlNode(c, idDirPattern, idFilePattern){
			modified = true
		}
	}

	return modified
}
