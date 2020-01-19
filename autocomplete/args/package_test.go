package args

//go:generate mockgen -package=args -destination=./filesystem_mock_test.go github.com/spoke-d/clui/autocomplete/args FileSystem
//go:generate mockgen -package=args -destination=./file_mock_test.go os FileInfo
