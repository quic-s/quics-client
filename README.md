# quics-client

quics-client is a client for the QUIC-S. It is **continuous file synchornization tool** based on the QUIC protocol. 


## Features

#### 1. App Controller
Start, Reboot, Stop management of the app. And also, user can check the status of the app. All the management is done by the command line made by **cobra** library.

#### 2. App Config
By using **viper** library, environment variables is managed. Mostly it is used for the configuration of the Root Directory Path for sync. And also, it is used for the configuration of the server IP address and port number.

#### 3. Connection Manager
* Server Connection : User can connect with access token and IP address, then dial to server to send first message with user information. 
* Root Directory Connection : User can select the directory that user want to manage. Choose Local Root Directory or Remote Root Directory that any folder user wants. And then, user can start to sync.

#### 4. File Manager
By using **fsnotify** library and goroutine, user can monitor the file changes in the directory. And then, user can send the file changes to the server. 


#### 5. File Sync
For make sync in real time, the Sync Metadata is saved in the local db, the **BadgerDB**. This is for the comparison of the file changes. And also, user can check the file changes in the local db.

#### 6. Extension
**HTTP/3** is used for the each cobra command, except the controller command. And also, the API is used for the communication between the client and server. By using the HTTP/3 API, developers can easily integrate the GUI or other functions.

## Getting Started

### 1. Docker
    
```
    docker run -it --rm -v /path/to/your/dir:/data chromato99/quics-client
```

### 2. Local Install

- 1. Download the latest release from the [releases page]()
- 2. Unpack the archive.
- 3. Run `./qic` in the unpacked directory.



### 3. Build from source

- 1. Install Go 1.21 or later.
- 2. Clone this repository.
- 3. Run `go build ./cmd` in the root of the repository.


## How to use

| Command | Options | Description |
| --- | --- | --- |
| `qic start` | `--hPort` | Start the app. |
| `qic reboot` | | Reboot the app. |
| `qic shutdown` | | Shutdown the app. |
| `qic --help` | | Show the help message. |
| `qic config show` | | Show the config. |
| `qic config server` | `--host` `--port` | Set the server IP address and port number. |
| `qic config root` | `--abspath` `--name` | Set the root directory path and name. |
| `qic config delete` | `--key` | Delete the config. |
| `qic connect server` | `--password` | Connect to the server with password. |
| `qic connect server` | `--host` `--port` `--password` | Connect to the server with IP address, port number and password. |
| `qic disconnect server` | `--password` | Disconnect to the server with password. |
| `qic connect root` | `--local` `--password` | Connect to the local root directory with password. |
| `qic connect root` | `--remote` `--password` | Connect to the remote root directory with password. |
| `qic connect list-remote` | | List the remote root directory. |
| `qic disconnect root` | `--root` `--password` | Disconnect to the root directory with password. |
| `qic sync status` | `--pick` | Show the status of the sync Metadata. |
| `qic sync rescan` | | Ask Server rescan the sync. |



## Contribute

- To report bugs or request features, please use the issue tracker. Before you do so, make sure you are running the latest version, and please do a quick search to see if the issue has already been reported.

- For more discussion, please join the [quics discord](https://discord.gg/HRtY7pNZz2)

