# elabox-system-tools
Elabox runtime. Provides system that handles communication and execution of dapps.


## VS CODE DEBUGGER SETTING
This necessary capturing ide related conditions. foundation.system.IDE constant is true if currently debugging inside ide.
- Open launch.json
- add buildflags property to launch configuration if not yet existed "buildFlags": "-tags=IDE"

## Building from Source
- Configure for building
```
cd ./internal/scripts
./setup.sh
```

- Building
```
cd ./internal/scripts
./build.sh
<follow sh instructions>
```

## Consepts
An operating system designed with holistic approach and security in mind powered by blockchains.
Features:
    - Decentralized Identity(DID) - access to the system and its data are protected by your identity.
    - Smart Contract 
    - Dapp Store - install and publish 3rd party dapps via dapp store 

Dapp Architecture
- a Dapp can talk to each other, broadcast and listen to any events
and handle background tasks. This is made possible via json RPC on the top of socket.
- Apps are being managed and handled by AppMan( Application Manager).
- a dapp should be defined on info.json
- a dapp compose of 3 main components.

Service - a  component that runs on background and provides services that can be consumed by other dapps.

Activity - a component that would only run when started by user or another dapp. Usually this component contain UI to handle user operations. 

Broadcast Listener - a component that would only executes after an event was triggered. Analogy
Sample Broadcast listener

