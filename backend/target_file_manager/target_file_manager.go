package target_file_manager

import "salami/common/types"

func VerifyChecksums(targetFilesMeta []*types.TargetFileMeta) error {
	// Verify checksums
	return nil
}

func UpdateTargetFilesMeta() ([]*types.TargetFileMeta, error) {
	// Update checksums
	return []*types.TargetFileMeta{}, nil
}

func WriteTargetFile(filePath string, content string) error {
	// Write target files
	return nil
}

// func (b *Backend) writeFiles(codeFiles []*types.DestinationFile) []error {
// 	var wg sync.WaitGroup
// 	errorsChan := make(chan error)

// 	for _, codeFile := range codeFiles {
// 		wg.Add(1)
// 		go func(cf *types.DestinationFile) {
// 			defer wg.Done()
// 			err := b.writeFile(*cf)
// 			if err != nil {
// 				errorsChan <- err
// 			}
// 		}(codeFile)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(errorsChan)
// 	}()
// 	errors := []error{}
// 	for err := range errorsChan {
// 		errors = append(errors, err)
// 	}
// 	if len(errors) > 0 {
// 		return errors
// 	}
// 	return nil
// }

// func (b *Backend) writeFile(codeFile types.DestinationFile) error {
// 	compilerConfig := config.GetCompilerConfig()
// 	fullPath := path.Join(compilerConfig.TargetDir, codeFile.FilePath)

// 	err := os.MkdirAll(path.Dir(fullPath), os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	if _, err := os.Stat(fullPath); err == nil {
// 		err = os.Remove(fullPath)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	file, err := os.Create(fullPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = file.WriteString(codeFile.Content)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
