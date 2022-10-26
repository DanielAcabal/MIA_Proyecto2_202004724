package estructuras
type Pares struct{
	Inicio int
	Tamanio int
}
type Particion struct{
    //Atributos
    Part_status [1]byte
    Part_type [1]byte
    Part_fit [2]byte
    Part_start [10]byte
    Part_size [10]byte
    Part_name [50]byte
}
type MBR struct{
	Mbr_tamano [10]byte
	Mbr_fecha_creacion [50]byte
	Mbr_dsk_signature [5]byte
	Dsk_fit [1]byte
	Mbr_partition_1 Particion
	Mbr_partition_2 Particion
	Mbr_partition_3 Particion
	Mbr_partition_4 Particion
}