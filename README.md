# elabox-system-tools
Tools to manage different parts of Elabox


## VS CODE DEBUGGER SETTING
This necessary capturing ide related conditions. foundation.system.IDE constant is true if currently debugging inside ide.
- Open laun.json
- add buildflags property to launch configuration if not yet existed "buildFlags": "-tags=IDE"