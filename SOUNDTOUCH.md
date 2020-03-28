
## NowPlaying Status

```golang
const (
	PLAY_STATE          = "PLAY_STATE"
	PAUSE_STATE         = "PAUSE_STATE"
	BUFFERING_STATE     = "BUFFERING_STATE"
	INVALID_PLAY_STATUS = "INVALID_PLAY_STATUS"
	STOP_STATE          = "STOP"
	STANDBY             = "STANDBY"
)
```

| Description | Source | Action | PlayState | SourceAccount | DeviceId | Track | Artist | Album | Art | StreamType |  C - Src | C - Type |  C - Name | C - Loc | C - IsPresetable |
| ----------- | ------ | ------ | --------- | :-------------: | :--------: | :-----: | :------:| :-----: | :---: | :-------: | :-------: | :-------: |  :--------: | :-------: | :----------------: |
| Playing INTERNET_RADIO | LOCAL_INTERNET_RADIO | playing | PLAY_STATE | - | *x* | *x* | - | - | - | RADIO_STREAMING | - | LOCAL_INTERNET_RADIO | *x* | *x* | _true_ | 
| Playing TUNEIN         | TUNEIN | playing | PLAY_STATE | - | *x* | *x* | *x* | - | *x* | RADIO_STREAMING | TUNEIN | _stationurl_ | *x* | *x* | _true_ | 
| Playing AUX | AUX | playing | PLAX_STATE | - | *x* | - | - | - | - | - |  AUX | - |  "AUX IN" | - | _true_ |
| Soundbar playing TV | PRODUCT | Playing | PLAY_STATE | TV | *x* | - | - | - | - | - |  PRODUCT | - |  - | - | _false_ |
| Soundbar searching BLUETOOTH device | BLUETOOTH | searching | INVALID_PLAY_STATUS | - | *x* | - | - | - | - | - |  BLUETOOTH | - |  - | - | _false_ |
| Playing BLUETOOTH | BLUETOOTH | playing | PLAY_STATE | - | *x* | *x* | *x* | *x* | - | - |  BLUETOOTH | - |  *x* | - | _false_ |
| Pausing BLUETOOTHG | BLUETOOTH | playing | PAUSE_STATE | - | *x* | *x* | *x* | *x* | - | - |  BLUETOOTH | - |  *x* | - | _false_ |
| Pausing a STORED_MUSIC | STORED_MUSIC | pausing | PAUSE_STATE | *x* | *x* | *x* | *x* | *x* | *x* | - | STORED_MUSIC | - | *x* | *x* | _true_ |
| Standby w STORED_MUSIC | STANDBY | standby | - | - | *x* | - | - | - | - | - |  STANBY | - |  - | - | - |
| Playing STORED_MUSIC | STORED_MUSIC | playing | PLAY_STATE | *x* | *x* | *x* | *x* | *x* | *x* | - |  STORED_MUSIC | -|  *x* | *x* | _true_ |
