LL_PATH=../.

if [ -d $LL_PATH ]; then
	cd $LL_PATH
	patch -p1 -R < .docker/gopacket_pcap.patch
fi
