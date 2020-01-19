package install

//go:generate mockgen -package=install -destination=./executable_mock_test.go github.com/spoke-d/clui/autocomplete/install Executable
//go:generate mockgen -package=install -destination=./filesystem_mock_test.go github.com/spoke-d/clui/autocomplete/fsys FileSystem,File
//go:generate mockgen -package=install -destination=./user_mock_test.go github.com/spoke-d/clui/autocomplete/install User
