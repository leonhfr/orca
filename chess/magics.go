package chess

// rookMagics contains the magics for rooks indexed by square.
//
// This literal has been automatically generated, do not edit.
var rookMagics = [64]Magic{
	{Mask: 0x101010101017E, Magic: 0x80008050C00024, Shift: 52, Offset: 0},
	{Mask: 0x202020202027C, Magic: 0x40009008200041, Shift: 53, Offset: 4096},
	{Mask: 0x404040404047A, Magic: 0x280100081A82000, Shift: 53, Offset: 6144},
	{Mask: 0x8080808080876, Magic: 0x80100080080204, Shift: 53, Offset: 8192},
	{Mask: 0x1010101010106E, Magic: 0x80240008001A80, Shift: 53, Offset: 10240},
	{Mask: 0x2020202020205E, Magic: 0x500240012450008, Shift: 53, Offset: 12288},
	{Mask: 0x4040404040403E, Magic: 0x4280008006002100, Shift: 53, Offset: 14336},
	{Mask: 0x8080808080807E, Magic: 0x680002140800100, Shift: 52, Offset: 16384},
	{Mask: 0x1010101017E00, Magic: 0x4000800080400820, Shift: 53, Offset: 20480},
	{Mask: 0x2020202027C00, Magic: 0x20400044A01008, Shift: 54, Offset: 22528},
	{Mask: 0x4040404047A00, Magic: 0x4D004020003300, Shift: 54, Offset: 23552},
	{Mask: 0x8080808087600, Magic: 0x4000801000809800, Shift: 54, Offset: 24576},
	{Mask: 0x10101010106E00, Magic: 0x201800800804400, Shift: 54, Offset: 25600},
	{Mask: 0x20202020205E00, Magic: 0x9001803000400, Shift: 54, Offset: 26624},
	{Mask: 0x40404040403E00, Magic: 0x405000100020004, Shift: 54, Offset: 27648},
	{Mask: 0x80808080807E00, Magic: 0x460800041000080, Shift: 53, Offset: 28672},
	{Mask: 0x10101017E0100, Magic: 0x41800A4002452001, Shift: 53, Offset: 30720},
	{Mask: 0x20202027C0200, Magic: 0x410808040002008, Shift: 54, Offset: 32768},
	{Mask: 0x40404047A0400, Magic: 0x800410020043100, Shift: 54, Offset: 33792},
	{Mask: 0x8080808760800, Magic: 0x440220008120240, Shift: 54, Offset: 34816},
	{Mask: 0x101010106E1000, Magic: 0xD02020020081024, Shift: 54, Offset: 35840},
	{Mask: 0x202020205E2000, Magic: 0x2808002000400, Shift: 54, Offset: 36864},
	{Mask: 0x404040403E4000, Magic: 0x410040012D00803, Shift: 54, Offset: 37888},
	{Mask: 0x808080807E8000, Magic: 0x200E0000640081, Shift: 53, Offset: 38912},
	{Mask: 0x101017E010100, Magic: 0x200A08480104000, Shift: 53, Offset: 40960},
	{Mask: 0x202027C020200, Magic: 0x900200080804000, Shift: 54, Offset: 43008},
	{Mask: 0x404047A040400, Magic: 0x800408600120220, Shift: 54, Offset: 44032},
	{Mask: 0x8080876080800, Magic: 0x401200200A01, Shift: 54, Offset: 45056},
	{Mask: 0x1010106E101000, Magic: 0x8008000A800C0080, Shift: 54, Offset: 46080},
	{Mask: 0x2020205E202000, Magic: 0x81181400800A0080, Shift: 54, Offset: 47104},
	{Mask: 0x4040403E404000, Magic: 0x8004104000826B0, Shift: 54, Offset: 48128},
	{Mask: 0x8080807E808000, Magic: 0x32138820000590C, Shift: 53, Offset: 49152},
	{Mask: 0x1017E01010100, Magic: 0xB0204003800580, Shift: 53, Offset: 51200},
	{Mask: 0x2027C02020200, Magic: 0x288101004000, Shift: 54, Offset: 53248},
	{Mask: 0x4047A04040400, Magic: 0x100504101002000, Shift: 54, Offset: 54272},
	{Mask: 0x8087608080800, Magic: 0x100101002018, Shift: 54, Offset: 55296},
	{Mask: 0x10106E10101000, Magic: 0x810400800800, Shift: 54, Offset: 56320},
	{Mask: 0x20205E20202000, Magic: 0x81040080800200, Shift: 54, Offset: 57344},
	{Mask: 0x40403E40404000, Magic: 0x5080C001006, Shift: 54, Offset: 58368},
	{Mask: 0x80807E80808000, Magic: 0x12001441020000A4, Shift: 53, Offset: 59392},
	{Mask: 0x17E0101010100, Magic: 0x4010400090668000, Shift: 53, Offset: 61440},
	{Mask: 0x27C0202020200, Magic: 0x30094020004000, Shift: 54, Offset: 63488},
	{Mask: 0x47A0404040400, Magic: 0x1060012641010010, Shift: 54, Offset: 64512},
	{Mask: 0x8760808080800, Magic: 0x2088001000210100, Shift: 54, Offset: 65536},
	{Mask: 0x106E1010101000, Magic: 0x850008010010, Shift: 54, Offset: 66560},
	{Mask: 0x205E2020202000, Magic: 0x4002010040128, Shift: 54, Offset: 67584},
	{Mask: 0x403E4040404000, Magic: 0x240A8122900C0008, Shift: 54, Offset: 68608},
	{Mask: 0x807E8080808000, Magic: 0x4000028104620014, Shift: 53, Offset: 69632},
	{Mask: 0x7E010101010100, Magic: 0x408000C0006000C0, Shift: 53, Offset: 71680},
	{Mask: 0x7C020202020200, Magic: 0x4810004002200040, Shift: 54, Offset: 73728},
	{Mask: 0x7A040404040400, Magic: 0x820043001842080, Shift: 54, Offset: 74752},
	{Mask: 0x76080808080800, Magic: 0xA402200102A00, Shift: 54, Offset: 75776},
	{Mask: 0x6E101010101000, Magic: 0x9200800800240180, Shift: 54, Offset: 76800},
	{Mask: 0x5E202020202000, Magic: 0x4020080240080, Shift: 54, Offset: 77824},
	{Mask: 0x3E404040404000, Magic: 0x400100906080400, Shift: 54, Offset: 78848},
	{Mask: 0x7E808080808000, Magic: 0x440204085240200, Shift: 53, Offset: 79872},
	{Mask: 0x7E01010101010100, Magic: 0x8420C680001301, Shift: 52, Offset: 81920},
	{Mask: 0x7C02020202020200, Magic: 0x200824022021302, Shift: 53, Offset: 86016},
	{Mask: 0x7A04040404040400, Magic: 0x200419C10011, Shift: 53, Offset: 88064},
	{Mask: 0x7608080808080800, Magic: 0x1000410000821, Shift: 53, Offset: 90112},
	{Mask: 0x6E10101010101000, Magic: 0x101002800100A15, Shift: 53, Offset: 92160},
	{Mask: 0x5E20202020202000, Magic: 0x80E000804301902, Shift: 53, Offset: 94208},
	{Mask: 0x3E40404040404000, Magic: 0x402A108021004, Shift: 53, Offset: 96256},
	{Mask: 0x7E80808080808000, Magic: 0x4000040049008022, Shift: 52, Offset: 98304},
}

// bishopMagics contains the magics for bishops indexed by square.
//
// This literal has been automatically generated, do not edit.
var bishopMagics = [64]Magic{
	{Mask: 0x40201008040200, Magic: 0x20230120410000, Shift: 58, Offset: 0},
	{Mask: 0x402010080400, Magic: 0x4140400740004, Shift: 59, Offset: 64},
	{Mask: 0x4020100A00, Magic: 0x49081103002400, Shift: 59, Offset: 96},
	{Mask: 0x40221400, Magic: 0x10040C8800200000, Shift: 59, Offset: 128},
	{Mask: 0x2442800, Magic: 0x4042007200100, Shift: 59, Offset: 160},
	{Mask: 0x204085000, Magic: 0x1100A00401000, Shift: 59, Offset: 192},
	{Mask: 0x20408102000, Magic: 0x1084038250008000, Shift: 59, Offset: 224},
	{Mask: 0x2040810204000, Magic: 0x1021010082A00000, Shift: 58, Offset: 256},
	{Mask: 0x20100804020000, Magic: 0x10902902044020, Shift: 59, Offset: 320},
	{Mask: 0x40201008040000, Magic: 0x50809180804C000, Shift: 59, Offset: 352},
	{Mask: 0x4020100A0000, Magic: 0x46010501440000, Shift: 59, Offset: 384},
	{Mask: 0x4022140000, Magic: 0x200020080000874, Shift: 59, Offset: 416},
	{Mask: 0x244280000, Magic: 0x88142C0C061000C0, Shift: 59, Offset: 448},
	{Mask: 0x20408500000, Magic: 0x8900012120000000, Shift: 59, Offset: 480},
	{Mask: 0x2040810200000, Magic: 0x52040088050010, Shift: 59, Offset: 512},
	{Mask: 0x4081020400000, Magic: 0x1500003202020014, Shift: 59, Offset: 544},
	{Mask: 0x10080402000200, Magic: 0x3D10100000, Shift: 59, Offset: 576},
	{Mask: 0x20100804000400, Magic: 0x602C0602240010, Shift: 59, Offset: 608},
	{Mask: 0x4020100A000A00, Magic: 0x2020801140C0000, Shift: 57, Offset: 640},
	{Mask: 0x402214001400, Magic: 0x200080802080022, Shift: 57, Offset: 768},
	{Mask: 0x24428002800, Magic: 0x4021403011000000, Shift: 57, Offset: 896},
	{Mask: 0x2040850005000, Magic: 0x40400040481A0080, Shift: 57, Offset: 1024},
	{Mask: 0x4081020002000, Magic: 0x4400081D04100000, Shift: 59, Offset: 1152},
	{Mask: 0x8102040004000, Magic: 0x2080005056065090, Shift: 59, Offset: 1184},
	{Mask: 0x8040200020400, Magic: 0x201109224040801, Shift: 59, Offset: 1216},
	{Mask: 0x10080400040800, Magic: 0x880010506002400B, Shift: 59, Offset: 1248},
	{Mask: 0x20100A000A1000, Magic: 0x100001004000, Shift: 57, Offset: 1280},
	{Mask: 0x40221400142200, Magic: 0x800001004E810400, Shift: 55, Offset: 1408},
	{Mask: 0x2442800284400, Magic: 0x14E80200044A0001, Shift: 55, Offset: 1920},
	{Mask: 0x4085000500800, Magic: 0x4D400C9820100, Shift: 57, Offset: 2432},
	{Mask: 0x8102000201000, Magic: 0x8910020013CA006C, Shift: 59, Offset: 2560},
	{Mask: 0x10204000402000, Magic: 0x8084011034000, Shift: 59, Offset: 2592},
	{Mask: 0x4020002040800, Magic: 0x44406200000, Shift: 59, Offset: 2624},
	{Mask: 0x8040004081000, Magic: 0x300C0401060000, Shift: 59, Offset: 2656},
	{Mask: 0x100A000A102000, Magic: 0x14810444000, Shift: 57, Offset: 2688},
	{Mask: 0x22140014224000, Magic: 0x18000600202000, Shift: 55, Offset: 2816},
	{Mask: 0x44280028440200, Magic: 0x200036020428800, Shift: 55, Offset: 3328},
	{Mask: 0x8500050080400, Magic: 0x3010601050020, Shift: 57, Offset: 3840},
	{Mask: 0x10200020100800, Magic: 0x100008400814033, Shift: 59, Offset: 3968},
	{Mask: 0x20400040201000, Magic: 0x6010008908008040, Shift: 59, Offset: 4000},
	{Mask: 0x2000204081000, Magic: 0x100500610012008, Shift: 59, Offset: 4032},
	{Mask: 0x4000408102000, Magic: 0x12000110A8400080, Shift: 59, Offset: 4064},
	{Mask: 0xA000A10204000, Magic: 0x4404402C008000, Shift: 57, Offset: 4096},
	{Mask: 0x14001422400000, Magic: 0x8148898144006004, Shift: 57, Offset: 4224},
	{Mask: 0x28002844020000, Magic: 0x1682808800004, Shift: 57, Offset: 4352},
	{Mask: 0x50005008040200, Magic: 0x8000840C4000001, Shift: 57, Offset: 4480},
	{Mask: 0x20002010080400, Magic: 0x90284200000, Shift: 59, Offset: 4608},
	{Mask: 0x40004020100800, Magic: 0x400008082800200, Shift: 59, Offset: 4640},
	{Mask: 0x20408102000, Magic: 0x1420040424400000, Shift: 59, Offset: 4672},
	{Mask: 0x40810204000, Magic: 0x400A405200020, Shift: 59, Offset: 4704},
	{Mask: 0xA1020400000, Magic: 0x8144A80208040D10, Shift: 59, Offset: 4736},
	{Mask: 0x142240000000, Magic: 0x80800822000, Shift: 59, Offset: 4768},
	{Mask: 0x284402000000, Magic: 0x8270A00004120030, Shift: 59, Offset: 4800},
	{Mask: 0x500804020000, Magic: 0x8880000439022008, Shift: 59, Offset: 4832},
	{Mask: 0x201008040200, Magic: 0x82208124880, Shift: 59, Offset: 4864},
	{Mask: 0x402010080400, Magic: 0x48114010808, Shift: 59, Offset: 4896},
	{Mask: 0x2040810204000, Magic: 0x44002404064800, Shift: 58, Offset: 4928},
	{Mask: 0x4081020400000, Magic: 0xC00028888091002, Shift: 59, Offset: 4992},
	{Mask: 0xA102040000000, Magic: 0x4014080480240, Shift: 59, Offset: 5024},
	{Mask: 0x14224000000000, Magic: 0x400088080120A0, Shift: 59, Offset: 5056},
	{Mask: 0x28440200000000, Magic: 0x8800800000200100, Shift: 59, Offset: 5088},
	{Mask: 0x50080402000000, Magic: 0x2000020208120400, Shift: 59, Offset: 5120},
	{Mask: 0x20100804020000, Magic: 0x8C008608080102, Shift: 59, Offset: 5152},
	{Mask: 0x40201008040200, Magic: 0x80200202434100, Shift: 58, Offset: 5184},
}