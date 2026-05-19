package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if !ValidUntil.IsZero() && time.Now().After(ValidUntil) {
		fmt.Println("Program has expired")
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Please drag and drop an APK file onto this program")
		fmt.Println("Or manually enter the APK file path")
		fmt.Print("Enter APK file path: ")
		var input string
		fmt.Scanln(&input)
		if input == "" {
			fmt.Println("No file path entered")
			fmt.Println("Press Enter to exit...")
			fmt.Scanln()
			return
		}
		os.Args = append(os.Args, input)
	}

	apkPath := os.Args[1]
	apkPath = strings.Trim(apkPath, "\"")
	apkPath = strings.TrimSpace(apkPath)

	if !strings.HasSuffix(strings.ToLower(apkPath), ".apk") {
		apkPath += ".apk"
	}

	if _, err := os.Stat(apkPath); os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", apkPath)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	tmpDir, err := os.MkdirTemp("", "apksigner_*")
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("Extracting tool files...")

	aaptPath := filepath.Join(tmpDir, "aapt.exe")
	apksignerPath := filepath.Join(tmpDir, "apksigner.jar")
	zipalignPath := filepath.Join(tmpDir, "zipalign.exe")
	keystorePath := filepath.Join(tmpDir, KeystoreFile)

	extractFile("aapt.exe", aaptPath)
	extractFile("apksigner.jar", apksignerPath)
	extractFile("zipalign.exe", zipalignPath)
	extractFile("key.jks", keystorePath)

	apkAbs, _ := filepath.Abs(apkPath)
	apkDir := filepath.Dir(apkAbs)
	apkBase := filepath.Base(apkAbs)
	apkName := strings.TrimSuffix(apkBase, filepath.Ext(apkBase))

	javaPath := "java"
	if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
		javaPath = filepath.Join(javaHome, "bin", "java")
	}

	fmt.Println("Getting APK information...")
	pkg, err := GetAPKPackage(apkAbs, aaptPath)
	if err != nil {
		fmt.Printf("Failed to get package name: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}
	fmt.Printf("Package name: %s\n", pkg)

	if !IsPackageAllowed(pkg) {
		fmt.Println("This package name is not in the allowed list for signing")
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	zipalignOutput := filepath.Join(apkDir, apkName+"_zipalign.apk")
	outputName := filepath.Join(apkDir, apkName+"_"+Platform+"_sign.apk")

	fmt.Println("Performing file alignment...")
	err = Zipalign(apkAbs, zipalignOutput, zipalignPath)
	if err != nil {
		fmt.Printf("File alignment failed: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	if _, err := os.Stat(zipalignOutput); os.IsNotExist(err) {
		fmt.Println("zipalign failed to generate aligned file")
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	fmt.Println("Performing signing...")
	err = SignAPK(zipalignOutput, outputName, keystorePath, javaPath, apksignerPath)
	if err != nil {
		fmt.Printf("Signing failed: %v\n", err)
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
		return
	}

	CleanupTempFiles(apkDir, apkName)

	fmt.Printf("Signing completed: %s\n", outputName)
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func extractFile(name string, dest string) {
	data, err := embeddedFiles.ReadFile(name)
	if err != nil {
		fmt.Printf("Failed to read embedded file %s: %v\n", name, err)
		return
	}
	err = os.WriteFile(dest, data, 0755)
	if err != nil {
		fmt.Printf("Failed to write file %s: %v\n", name, err)
	}
}
