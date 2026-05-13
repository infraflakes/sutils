package cd

import "fmt"

func GenerateInit(shellName string) (string, error) {
	switch shellName {
	case "fish":
		return fmt.Sprintf(`function scd
    set -l target (sn cd $argv)
    if test $status -eq 0; and test -d "$target"
        builtin cd "$target"
    else if test -n "$target"
        printf "%%s\n" $target
    end
end`), nil

	case "zsh":
		return `scd() {
    local target
    target=$(sn cd "$@")
    if [ $? -eq 0 ] && [ -d "$target" ]; then
        builtin cd "$target"
    else
        [ -n "$target" ] && echo "$target"
    fi
}`, nil

	case "bash":
		return `scd() {
    local target
    target=$(sn cd "$@")
    if [ $? -eq 0 ] && [ -d "$target" ] ; then
        builtin cd "$target"
    else
        [ -n "$target" ] && echo "$target"
    fi
}`, nil

	default:
		return "", fmt.Errorf("unsupported shell: %s. Supported shells: fish, bash, zsh", shellName)
	}
}
