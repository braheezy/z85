Require echo
Require z85
Require hexyl

Set Shell "bash"
Set FontSize 24
Set WindowBar ColorfulRight
Set Width 1350
Set Height 525
Set Theme "Catppuccin Macchiato"
Set Margin 30
Set MarginFill "#c6a0f6"
Set BorderRadius 10
Set LineHeight 1.75
Set Padding 10

Type "z85 -mode encode 864FD26FB559F75B"
Enter
Sleep 500ms

Type "echo 'HelloWorld' | z85 -mode decode | hexyl -P --no-characters --panels 1"
Enter

Type "z85 # Run with no args for full usage"
Sleep 500ms

Screenshot assets/demo.png
Sleep 500ms
