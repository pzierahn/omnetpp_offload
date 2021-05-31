# Set up raspberry as a worker

```shell
sudo timedatectl set-timezone Europe/Berlin

# https://github.com/lucas-clemente/quic-go/wiki/UDP-Receive-Buffer-Size
sudo sysctl -w net.core.rmem_max=2500000
```