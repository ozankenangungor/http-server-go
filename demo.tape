Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1100
Set Height 650
Set Theme "Dracula"
Set TypingSpeed 40ms
Set Padding 20

Type "./your_program.sh --directory /tmp/demo-dir > /tmp/demo-server.log 2>&1 &"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms

Type "curl -s http://localhost:4221/echo/hello-world -w '\n'"
Enter
Sleep 1800ms

Type `curl -s -H "User-Agent: curl-demo/1.0" http://localhost:4221/user-agent -w '\n'`
Enter
Sleep 1800ms

Type "curl -s http://localhost:4221/files/hello.txt -w '\n'"
Enter
Sleep 1800ms

Type `curl -si -H "Accept-Encoding: gzip" http://localhost:4221/echo/compressed | head -4`
Enter
Sleep 2500ms
