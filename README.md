This tool samples the default audio input device, 
performs an FFT and sends 32 bytes of spectrum analyzer band heights (0..255) 
to [nyukomatic](https://github.com/alexanderk23/nyukomatic/) via WebSocket 
using the [patched Bonzomatic server](https://github.com/alexanderk23/BonzomaticServer) protocol.

See [example.asm](example.asm) for an example of usage.

## Usage

```sh
./nm-fft-example2 -port 9000 -room test -r 48000
```
