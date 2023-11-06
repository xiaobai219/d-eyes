import "pe"

rule KelihosHlux
{
meta:
	description = "Detect the risk of Botnet Malware Kelihos Rule 1"
  strings:
    $KelihosHlux_HexString = { 73 20 7D 8B FE 95 E4 12 4F 3F 99 3F 6E C8 28 26 C2 41 D9 8F C1 6A 72 A6 CE 36 0F 73 DD 2A 72 B0 CC D1 07 8B 2B 98 73 0E 7E 8C 07 DC 6C 71 63 F4 23 27 DD 17 56 AE AB 1E 30 52 E7 54 51 F7 20 ED C7 2D 4B 72 E0 77 8E B4 D2 A8 0D 8D 6A 64 F9 B7 7B 08 70 8D EF F3 9A 77 F6 0D 88 3A 8F BB C8 89 F5 F8 39 36 BA 0E CB 38 40 BF 39 73 F4 01 DC C1 17 BF C1 76 F6 84 8F BD 87 76 BC 7F 85 41 81 BD C6 3F BC 39 BD C0 89 47 3E 92 BD 80 60 9D 89 15 6A C6 B9 89 37 C4 FF 00 3D 45 38 09 CD 29 00 90 BB B6 38 FD 28 9C 01 39 0E F9 30 A9 66 6B 19 C9 F8 4C 3E B1 C7 CB 1B C9 3A 87 3E 8E 74 E7 71 D1 }
  condition:
    $KelihosHlux_HexString
}

rule kelihos_botnet_pdb {
meta:
	description = "Detect the risk of Botnet Malware Kelihos Rule 2"
    hash = "f0a6d09b5f6dbe93a4cf02e120a846073da2afb09604b7c9c12b2e162dfe7090"
strings:
	$pdb = "\\Only\\Must\\Not\\And.pdb"
	$pdb1 = "\\To\\Access\\Do.pdb"
condition:
	uint16(0) == 0x5a4d and
	filesize < 1440KB and
	any of them
}