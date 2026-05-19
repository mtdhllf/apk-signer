package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Zipalign(inputPath string, outputPath string, zipalignPath string) error {
	cmd := exec.Command(zipalignPath, "-p", "-f", "-v", "4", inputPath, outputPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("zipalign failed: %v", err)
	}
	return nil
}

func SignAPK(inputPath string, outputPath string, keystorePath string, javaPath string, apksignerPath string) error {
	absKeystore, _ := filepath.Abs(keystorePath)
	absInput, _ := filepath.Abs(inputPath)
	absOutput, _ := filepath.Abs(outputPath)

	cmd := exec.Command(javaPath, "-jar", apksignerPath,
		"sign",
		"--ks", absKeystore,
		"--ks-key-alias", KeystoreAlias,
		"--ks-pass", "pass:"+StorePassword,
		"--key-pass", "pass:"+KeyPassword,
		"--out", absOutput,
		absInput)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("signing failed: %v", err)
	}
	return nil
}

func CleanupTempFiles(dir string, baseName string) {
	zipalignFile := filepath.Join(dir, baseName+"_zipalign.apk")
	idsigFile := filepath.Join(dir, baseName+"_"+Platform+"_sign.apk.idsig")

	os.Remove(zipalignFile)
	os.Remove(idsigFile)
}
