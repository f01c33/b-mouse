# binary search mouse control (b-mouse)

Allows to control your mouse with your keyboard, using binary search in your monitor's 2d space for cursor control

![demo.gif](./demo.gif)

#### install

You need to install [robotgo's](https://github.com/robotn/robotgo) dependencies and [glfw's](https://github.com/go-gl/glfw) dependencies, then you can do the next step:

```bash
go install github.com/f01c33/b-mouse@latest
```

#### usage

you need $GOPATH/bin (usually ``` ~/go/bin ```) on your $PATH (PATH environment variable on windows)

```bash
b-mouse #on your terminal
```

You can use Vim bindings, the arrow keys or wasd to control the mouse

enter or 1 for left mouse click

shift or 2 for middle mouse click

space or 3 for left mouse click

pause to hide/show the program

esc or q to reset the position to the start position