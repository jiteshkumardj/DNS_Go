# DNS Catcher

The DNS catcher is intended to intercept DNS packets and print their information to the terminal.

## Dependencies

DNA Catcher requires a `libpcap` library.

### Linux

```shell
sudo apt-get install libpcap-dev
```

### Windows
Windows requires `winpcap` or `npcap`. If both are installed at the same time, `npcap` is preferred. Both can be installed using windows installer:

`winpcap` - https://www.winpcap.org/install/default.htm
`npcap` - https://npcap.com/#download
