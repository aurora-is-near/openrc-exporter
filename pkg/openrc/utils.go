package openrc

// #cgo LDFLAGS: -lrc
/*
#include <rc.h>
#include <stdlib.h>
*/
import "C"

// struct rc_stringlist {
//     struct rc_string *tqh_first;
//     struct rc_string **tqh_last;
// }
//
// struct rc_string {
// 	   char *value;
// 	   struct {
//		  rc_string *tqe_next;
//		  rc_string *tqe_prev
//	   } entries;
// }
func goStringList(stringlist *C.struct_rc_stringlist) []string {
	if stringlist == nil {
		return []string{}
	}

	item := (*stringlist).tqh_first
	last := *((*stringlist).tqh_last)

	var list []string

	for item != last {
		value := (*item).value
		list = append(list, C.GoString(value))
		item = item.entries.tqe_next
	}

	return list
}
