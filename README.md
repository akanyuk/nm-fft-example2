This tool samples the default audio input device, 
performs an FFT and sends 32 bytes of spectrum analyzer band heights (0..255) 
to [nyukomatic](https://github.com/alexanderk23/nyukomatic/) via WebSocket 
using the [patched Bonzomatic server](https://github.com/alexanderk23/BonzomaticServer) protocol.

## Usage

1. Run this program. If the default port is used, you can change it.
```sh
./nm-fft-example2 -port 9000 -r 48000
``` 
2. Launch the [nyukomatic](https://github.com/alexanderk23/nyukomatic) and set the server URL: `ws://127.0.0.1:9000/test/nyuk` (press F1).
3. Copy/paste [example.asm](example.asm) into the editor.
4. Play music, adjust the default microphone volume.