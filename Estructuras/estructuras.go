package estructuras

import (
	"strconv"
	"time"
	"unsafe"
)
type Pares struct{
	Inicio int64
	Tamanio int64
}
type Particion struct{
    //Atributos
    Part_status [1]byte
    Part_type [1]byte
    Part_fit [2]byte
    Part_start [5]byte
    Part_size [5]byte
    Part_name [15]byte
}
type MBR struct{
	Mbr_tamano [8]byte
	Mbr_fecha_creacion [15]byte
	Mbr_dsk_signature [5]byte
	Dsk_fit [1]byte
	Mbr_partition_1 Particion
	Mbr_partition_2 Particion
	Mbr_partition_3 Particion
	Mbr_partition_4 Particion
}
type Pmontada struct{
    Id string;
    Particion Particion;
    Path string;
	TiempoM string;
}
type SuperBloque struct{
    S_filesystem_type [1]byte
    S_inodes_count [4]byte    //Int->[]
    S_blocks_count [4]byte   //Int->[]
    S_free_blocks_count [4]byte
    S_free_inodes_count [4]byte
    S_mtime [30]byte
    S_mnt_count [1]byte
    S_magic [4]byte
    S_inode_size [4]byte
    S_block_size [4]byte
    S_firts_ino [4]byte
    S_first_blo [4]byte
    S_bm_inode_start [4]byte  //Into->[]
    S_bm_block_start [4]byte //Int->[]
    S_inode_start [4]byte
    S_block_start [4]byte
}
type Inodo struct{
    I_uid [2]byte
    I_gid [2]byte
    I_size [4]byte
    I_atime [8]byte
    I_ctime [8]byte
    I_mtime [8]byte
    I_block [16]byte
    I_type [1]byte
    I_perm [3]byte
}
type content struct{
    B_name [12]byte
    B_inodo [4]byte
}
type BloqueCarpeta struct{ //tamanio es de 64 bytes
    B_content [4]content
}
type BloqueArchivos struct{// también es de tamanio 64 bytes
    B_content [64]byte
}
func NuevoInodo(nuevo *Inodo,uId int64,gId int64,size int64,tipo string,perm string){
    copy(nuevo.I_uid[:],IntToByteArray(uId))
    copy(nuevo.I_gid[:],IntToByteArray(gId))
    copy(nuevo.I_size[:],IntToByteArray(size))
    creado := time.Now()
    copy(nuevo.I_atime[:],creado.String())
    copy(nuevo.I_ctime[:],creado.String())
    copy(nuevo.I_mtime[:],creado.String())
    nulos := make([]byte,16)
    for i := 0; i < len(nulos); i++ {
        nulos[i] = IntToByteArray(-1)[0]    }
    copy(nuevo.I_block[:],nulos)
    
    copy(nuevo.I_type[:],tipo)
    copy(nuevo.I_perm[:],perm)
}
func NuevoBloqueCarpeta(nuevo *BloqueCarpeta){
    for i := 0; i < len(nuevo.B_content); i++ {
        copy(nuevo.B_content[i].B_name[:],"")
        copy(nuevo.B_content[i].B_inodo[:],IntToByteArray(-1)) //modificar    
    }
    
}
func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	cos := strconv.FormatInt(num, 36)
	copy(arr,cos)
	return arr
}
type Usuario struct{
    Id string
    User string
    Password string
}