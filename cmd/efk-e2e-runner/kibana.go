package main

// RunKibanaE2ETest is a func to run e2e tests for Kibana using API.
func RunKibanaE2ETest() error {
	err := createKibanaIndexPattern()
	if err != nil {
		return err
	}

	err = deleteKibanaIndexPattern()
	if err != nil {
		return err
	}
	return nil
}

func createKibanaIndexPattern() error {
	return nil
}

func deleteKibanaIndexPattern() error {
	return nil
}
