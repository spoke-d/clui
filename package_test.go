package clui_test

//go:generate mockgen -package=mocks -destination=./mocks/autocomplete.go github.com/spoke-d/clui AutoCompleteInstaller
//go:generate mockgen -package=mocks -destination=./mocks/ui.go github.com/spoke-d/clui UI
//go:generate mockgen -package=mocks -destination=./mocks/command.go github.com/spoke-d/clui Command
