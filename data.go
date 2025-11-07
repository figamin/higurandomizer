package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SpriteInfo represents a sprite's expressions, folder, and preferred variant
type SpriteInfo struct {
	Expressions []string
	Folder      string
	Variant     string
}

// Prefix → folder lookup table
var FolderMap = map[string]string{
	"chibimion": "mion",
	"me":        "mion",
	"oisi":      "ooishi",
	"re":        "rena",
	"ri":        "rika",
	"sa":        "satoko",
	"ta":        "takano",
	"tie":       "chie",
	"tomi":      "tomitake",
	"kasa":      "kasai",
	"si":        "shion",
	"iri":       "irie",
	"aka":       "akane",
	"kei":       "keiichi",
	"sato":      "satoshi",
	"tetu":      "teppei",
	"rina":      "youhei",
	"aks":       "akasaka",
	"ha":        "hanyuu",
	"oko":       "fuko",
	"kameda":    "haruhi",
	"mo":        "eua",
	"mura":      "eua",
	"tamura":    "tamurahime",
	"une":       "une",
}
const (
    fuan_blush_close   = "fuan_blush_close"
    fuan_blush_open    = "fuan_blush_open"
    fuan_close         = "fuan_close"
    fuan_open          = "fuan_open"
    futeki_blush_close = "futeki_blush_close"
    futeki_blush_open  = "futeki_blush_open"
    futeki_close       = "futeki_close"
    futeki_open        = "futeki_open"
    L5_blush_close     = "L5_blush_close"
    L5_blush_open      = "L5_blush_open"
    L5_close           = "L5_close"
    L5_open            = "L5_open"
    normal_blush_close = "normal_blush_close"
    normal_blush_open  = "normal_blush_open"
    normal_close       = "normal_close"
    normal_open        = "normal_open"
    odoroki_blush_close = "odoroki_blush_close"
    odoroki_blush_open  = "odoroki_blush_open"
    odoroki_close       = "odoroki_close"
    odoroki_open        = "odoroki_open"
    sinken_blush_close  = "sinken_blush_close"
    sinken_blush_open   = "sinken_blush_open"
    sinken_close        = "sinken_close"
    sinken_open         = "sinken_open"
    smile_blush_close   = "smile_blush_close"
    smile_blush_open    = "smile_blush_open"
    smile_close         = "smile_close"
    smile_open          = "smile_open"
)

func generateVariants(n int) []string {
    v := make([]string, n)
    for i := 0; i < n; i++ {
        v[i] = fmt.Sprintf("v%03d", i+1)
    }
    return v
}

var spriteSets = generateVariants(55)

var PrefixFolderMappings = map[string]string{
    "chibimion": "mion",
    "tamura":    "tamurahime",
    "kameda":    "haruhi",
    "tomi":      "tomitake",
    "mura":      "eua",
    "sato":      "satoshi",
    "kasa":      "kasai",
    "oisi":      "ooishi",
    "iri":       "irie",
    "tetu":      "teppei",
    "rina":      "youhei",
    "aks":       "akasaka",
    "kei":       "keiichi",
    "aka":       "akane",
    "tie":       "chie",
    "oko":       "fuko",
    "ha":        "hanyuu",
    "me":        "mion",
    "re":        "rena",
    "ri":        "rika",
    "sa":        "satoko",
    "ta":        "takano",
    "si":        "shion",
    "mo":        "eua",
    "une":       "une",
}


// map mangagamer sprites to mei sprites
var RawGameSprites = map[string][]string{
	// mion 1a - normal school
	// ep1
	"me1a_akuwarai_a1_1": {futeki_open, spriteSets[0]},
	"me1a_akuwarai_a1_2": {futeki_open, spriteSets[0]},
	"me1a_def_a1_0": {smile_open, spriteSets[0]},
	"me1a_hig_maji_a1_0": {L5_open, spriteSets[0]},
	"me1a_huteki_a1_1": {futeki_open, spriteSets[0]},
	"me1a_huteki_a1_2": {futeki_open, spriteSets[0]},
	"me1a_ikari_a1_1": {sinken_open, spriteSets[0]},
	"me1a_majime_a1_0": {sinken_open, spriteSets[0]},
	"me1a_majime_a1_1": {sinken_open, spriteSets[0]},
	"me1a_odoroki_a1_1": {odoroki_open, spriteSets[0]},
	"me1a_tohoho_a1_0": {normal_open, spriteSets[0]},
	"me1a_tohoho_a1_1": {normal_open, spriteSets[0]},
	"me1a_tokui_a1_1": {futeki_close, spriteSets[0]},
	"me1a_tokui_a1_2": {futeki_close, spriteSets[0]},
	"me1a_warai_a1_1": {smile_close, spriteSets[0]},
	"me1a_warai_a1_2": {smile_close, spriteSets[0]},
	"me1a_wink_a1_1": {smile_close, spriteSets[0]},
	"me1a_wink_a1_2": {smile_close, spriteSets[0]},
	"me1a_yowaki_a1_1": {fuan_open, spriteSets[0]},
	"me1a_yowaki_a1_2": {fuan_open, spriteSets[0]},
	// ep2
	"me1a_hau_a1_1": {fuan_blush_open, spriteSets[0]},
	"me1a_odoroki_a1_2": {odoroki_open, spriteSets[0]},
	// ep3
	"me1a_ikari_a1_2": {sinken_open, spriteSets[0]},
	"me1a_sinmyou_a1_0": {smile_blush_open, spriteSets[0]},
	"me1a_sinmyou_a1_1": {smile_blush_open, spriteSets[0]},
	// ep4
	"me1a_odoroki_a1_0": {odoroki_open, spriteSets[0]},
	// ep5
	"me1a_hau_a1_0": {fuan_blush_open, spriteSets[0]},
	"me1a_yowaki_a1_0": {fuan_open, spriteSets[0]},
	// ep6
	"me1a_akuwarai_a1_0": {futeki_open, spriteSets[0]},
	"me1a_huteki_a1_0": {futeki_open, spriteSets[0]},
	"me1a_tokui_a1_0": {futeki_close, spriteSets[0]},
	"me1a_warai_a1_0": {smile_close, spriteSets[0]},
	"me1a_wink_a1_0": {smile_close, spriteSets[0]},
	// ep7
	"me1a_def_a1_1": {smile_open, spriteSets[0]},

	// mion 1b - thumbs school
	// ep1
	"me1b_akuwarai_a1_1": {futeki_open, spriteSets[0]},
	"me1b_akuwarai_a1_2": {futeki_open, spriteSets[0]},
	"me1b_def_a1_0": {smile_open, spriteSets[0]},
	"me1b_huteki_a1_1": {futeki_open, spriteSets[0]},
	"me1b_huteki_a1_2": {futeki_open, spriteSets[0]},
	"me1b_ikari_a1_1": {sinken_open, spriteSets[0]},
	"me1b_ikari_a1_2": {sinken_open, spriteSets[0]},
	"me1b_majime_a1_0": {sinken_open, spriteSets[0]},
	"me1b_odoroki_a1_1": {odoroki_open, spriteSets[0]},
	"me1b_odoroki_a1_2": {odoroki_open, spriteSets[0]},
	"me1b_tohoho_a1_0": {normal_open, spriteSets[0]},
	"me1b_tohoho_a1_1": {normal_open, spriteSets[0]},
	"me1b_tokui_a1_1": {futeki_close, spriteSets[0]},
	"me1b_tokui_a1_2": {futeki_close, spriteSets[0]},
	"me1b_warai_a1_1": {smile_close, spriteSets[0]},
	"me1b_wink_a1_1": {smile_close, spriteSets[0]},
	"me1b_wink_a1_2": {smile_close, spriteSets[0]},
	"me1b_yowaki_a1_1": {fuan_open, spriteSets[0]},
	"me1b_yowaki_a1_2": {fuan_open, spriteSets[0]},
	// ep2
	"me1b_hau_a1_1": {fuan_blush_open, spriteSets[0]},
	"me1b_majime_a1_1": {sinken_open, spriteSets[0]},
	"me1b_warai_a1_2": {smile_close, spriteSets[0]},
	// ep3
	"me1b_sinmyou_a1_0": {smile_blush_open, spriteSets[0]},
	"me1b_sinmyou_a1_1": {smile_blush_open, spriteSets[0]},
	// ep6
	"me1b_akuwarai_a1_0": {futeki_open, spriteSets[0]},
	"me1b_hau_a1_0": {futeki_open, spriteSets[0]},
	"me1b_odoroki_a1_0": {odoroki_open, spriteSets[0]},
	"me1b_tokui_a1_0": {futeki_close, spriteSets[0]},
	"me1b_warai_a1_0": {smile_close, spriteSets[0]},
	"me1b_wink_a1_0": {smile_close, spriteSets[0]},
	// ep7
	"me1b_def_a1_1": {smile_open, spriteSets[0]},

	// mion 2 - casual
	// ep1
	"me2_akuwarai_a1_1": {futeki_open, spriteSets[1]},
	"me2_akuwarai_a1_2": {futeki_open, spriteSets[1]},
	"me2_def_a1_0": {smile_open, spriteSets[1]},
	"me2_hig_maji_a1_0": {L5_open, spriteSets[1]},
	"me2_huteki_a1_1": {futeki_open, spriteSets[1]},
	"me2_huteki_a1_2": {futeki_open, spriteSets[1]},
	"me2_ikari_a1_1": {sinken_open, spriteSets[1]},
	"me2_ikari_a1_2": {sinken_open, spriteSets[1]},
	"me2_majime_a1_0": {sinken_open, spriteSets[1]},
	"me2_odoroki_a1_1": {odoroki_open, spriteSets[1]},
	"me2_tohoho_a1_0": {normal_open, spriteSets[1]},
	"me2_tohoho_a1_1": {normal_open, spriteSets[1]},
	"me2_tokui_a1_1": {futeki_close, spriteSets[1]},
	"me2_tokui_a1_2": {futeki_close, spriteSets[1]},
	"me2_warai_a1_1": {smile_close, spriteSets[1]},
	"me2_warai_a1_2": {smile_close, spriteSets[1]},
	"me2_wink_a1_1": {smile_close, spriteSets[1]},
	"me2_wink_a1_2": {smile_close, spriteSets[1]},
	// ep2
	"me2_def_a1_1": {smile_open, spriteSets[1]},
	"me2_hau_a1_1": {fuan_blush_open, spriteSets[1]},
	"me2_odoroki_a1_2": {odoroki_open, spriteSets[1]},
	"me2_sinmyou_a1_0": {smile_blush_open, spriteSets[1]},
	"me2_sinmyou_a1_1": {smile_blush_open, spriteSets[1]},
	"me2_yowaki_a1_1": {fuan_open, spriteSets[1]},
	// ep3
	"me2_yowaki_a1_2": {fuan_open, spriteSets[1]},
	// ep5
	"me2_akuwarai_a1_0": {futeki_open, spriteSets[1]},
	"me2_def_a1_2": {smile_open, spriteSets[1]},
	// ep6
	"me2_huteki_a1_0": {futeki_open, spriteSets[1]},
	"me2_odoroki_a1_0": {odoroki_open, spriteSets[1]},
	"me2_warai_a1_0": {smile_close, spriteSets[1]},
	"me2_wink_a1_0": {smile_close, spriteSets[1]},
	"me2_yowaki_a1_0": {fuan_open, spriteSets[1]},
	// ep7
	"me2_hau_a1_2": {fuan_blush_open, spriteSets[1]},

	// mion 3 - gym
	// ep1
	"me3_akuwarai_a1_1": {futeki_open, spriteSets[2]},
	"me3_akuwarai_a1_2": {futeki_open, spriteSets[2]},
	"me3_def_a1_0": {smile_open, spriteSets[2]},
	"me3_huteki_a1_2": {futeki_open, spriteSets[2]},
	"me3_tohoho_a1_0": {normal_open, spriteSets[2]},
	"me3_tokui_a1_1": {futeki_close, spriteSets[2]},
	"me3_tokui_a1_2": {futeki_close, spriteSets[2]},
	"me3_warai_a1_2": {smile_close, spriteSets[2]},
	"me3_wink_a1_1": {smile_close, spriteSets[2]},
	"me3_wink_a1_2": {smile_close, spriteSets[2]},
	// ep2
	"me3_huteki_a1_1": {futeki_open, spriteSets[2]},
	// ep6
	"me3_akuwarai_a1_0": {futeki_open, spriteSets[2]},
	"me3_huteki_a1_0": {futeki_open, spriteSets[2]},
	"me3_ikari_a1_0": {sinken_open, spriteSets[2]},
	"me3_majime_a1_0": {sinken_open, spriteSets[2]},
	"me3_odoroki_a1_0": {odoroki_open},
	"me3_tokui_a1_0": {futeki_close, spriteSets[2]},
	"me3_warai_a1_0": {smile_close, spriteSets[2]},
	"me3_wink_a1_0": {smile_close, spriteSets[2]},

	// mion 4 - school swimsuit
	// ep1
	"me4_akuwarai_a1_1": {futeki_open, spriteSets[5]},
	"me4_akuwarai_a1_2": {futeki_open, spriteSets[5]},
	"me4_huteki_a1_2": {futeki_open, spriteSets[5]},
	"me4_wink_a1_1": {smile_close, spriteSets[5]},
	"me4_wink_a1_2": {smile_close, spriteSets[5]},
	// rei
	"me4_tohoho_a1_0": {normal_open, spriteSets[5]},
	// hou+
	"me4_def_a1_1": {smile_open, spriteSets[5]},
	"me4_yowaki_a1_2": {fuan_open, spriteSets[5]},

	// mion 5 - punishment
	// ep6
	"me5_akuwarai_a1_0": {futeki_open, spriteSets[12]},
	"me5_def_a1_0": {smile_open, spriteSets[12]},
	"me5_huteki_a1_0": {futeki_open, spriteSets[12]},
	"me5_tohoho_a1_0": {normal_open, spriteSets[12]},
	"me5_warai_a1_0": {smile_close, spriteSets[12]},
	"me5_wink_a1_0": {smile_close, spriteSets[12]},

	// mion 7 - casual injured
	// hou+
	"me7_akuwarai_a1_2": {futeki_open, spriteSets[1]},
	"me7_def_a1_1": {smile_open, spriteSets[1]},
	"me7_huteki_a1_1": {futeki_open, spriteSets[1]},
	"me7_ikari_a1_2": {sinken_open, spriteSets[1]},
	"me7_majime_a1_0": {sinken_open, spriteSets[1]},
	"me7_odoroki_a1_2": {odoroki_open, spriteSets[1]},
	"me7_sinmyou_a1_0": {smile_blush_open, spriteSets[1]},
	"me7_tohoho_a1_0": {normal_open, spriteSets[1]},
	"me7_tokui_a1_2": {futeki_close, spriteSets[1]},
	"me7_warai_a1_2": {smile_close, spriteSets[1]},
	"me7_wink_a1_1": {smile_close, spriteSets[1]},
	"me7_yowaki_a1_2": {fuan_open, spriteSets[1]},

	// mion 8 - casual swimsuit
	// hou+
	"me8_akuwarai_a1_2": {futeki_open, spriteSets[5]},
	"me8_def_a1_1": {smile_open, spriteSets[5]},
	"me8_hau_a1_1": {fuan_blush_open, spriteSets[5]},
	"me8_huteki_a1_1": {futeki_open, spriteSets[5]},
	"me8_odoroki_a1_2": {odoroki_open, spriteSets[5]},
	"me8_tokui_a1_2": {futeki_close, spriteSets[5]},
	"me8_warai_a1_2": {smile_close, spriteSets[5]},
	"me8_wink_a1_1": {smile_close, spriteSets[5]},
	"me8_yowaki_a1_2": {smile_close, spriteSets[5]},

	// ooishi 1
	// ep1
	"oisi1_1_0": {smile_open, spriteSets[0]},
	"oisi1_2_0": {futeki_close, spriteSets[0]},
	"oisi1_2_1": {futeki_close, spriteSets[0]},
	"oisi1_3_2": {sinken_open, spriteSets[0]},
	// ep2
	"oisi1_5_0": {futeki_open, spriteSets[0]},
	"oisi1_5_2": {futeki_open, spriteSets[0]},
	// ep3
	"oisi1_4_1": {smile_open, spriteSets[0]},
	"oisi1_5_1": {futeki_open, spriteSets[0]},
	// ep4
	"oisi2_6_0": {smile_open, spriteSets[1]},
	"oisi2_7_0": {futeki_close, spriteSets[1]},
	"oisi2_8_2": {sinken_open, spriteSets[1]},
	"oisi2_9_1": {smile_open, spriteSets[1]},
	// ep5
	"oisi1_4_2": {smile_open, spriteSets[0]},
	// ep6
	"oisi1_3_0": {sinken_open, spriteSets[0]},
	"oisi1_4_0": {smile_open, spriteSets[0]},
	// ep7
	"oisi1_2_2": {futeki_close, spriteSets[0]},
	// rei
	"oisi2_7_2": {futeki_close, spriteSets[1]},
	"oisi2_10_2": {futeki_open, spriteSets[1]},

	// rena 1a - normal school
	// ep1
	"re1a_bikkuri_a1_2": {odoroki_blush_open, spriteSets[0]},
	"re1a_def_a1_0": {smile_blush_open, spriteSets[0]},
	"re1a_def_a1_2": {smile_blush_open, spriteSets[0]},
	"re1a_hau_a1_1": {fuan_blush_open, spriteSets[0]},
	"re1a_hig_def_a1_0": {L5_blush_open, spriteSets[0]},
	"re1a_hig_muhyou_a1_0": {L5_blush_open, spriteSets[0]},
	"re1a_kaii_a1_2": {smile_blush_close, spriteSets[0]},
	"re1a_komaru_a1_0": {fuan_blush_open, spriteSets[0]},
	"re1a_komaru_a2_0": {fuan_blush_open, spriteSets[0]},
	"re1a_nande_a1_1": {odoroki_blush_open, spriteSets[0]},
	"re1a_okoru_a1_0": {sinken_blush_open, spriteSets[0]},
	"re1a_okoru_a1_2": {sinken_blush_open, spriteSets[0]},
	"re1a_warai_a1_2": {smile_blush_close, spriteSets[0]},
	// ep4
	"re1a_nande_a1_0": {odoroki_blush_open, spriteSets[0]},
	// ep5
	"re1a_bikkuri_a1_1": {odoroki_blush_open, spriteSets[0]},
	// ep6
	"re1a_bikkuri_a1_0": {odoroki_blush_open, spriteSets[0]},
	"re1a_hau_a1_0": {fuan_blush_open, spriteSets[0]},
	"re1a_warai_a1_0": {smile_blush_close, spriteSets[0]},
	// ep7
	"re1a_hau_a1_2": {fuan_blush_open, spriteSets[0]},
	"re1a_hig_okoru_a1_2": {sinken_blush_open, spriteSets[0]},
	"re1a_nande_a1_2": {odoroki_blush_open, spriteSets[0]},
	

	// rena 1b - hands school
	// ep1
	"re1b_bikkuri_b1_2": {odoroki_blush_open, spriteSets[0]},
	"re1b_def_b1_0": {smile_blush_open, spriteSets[0]},
	"re1b_hau_b1_1": {fuan_blush_open, spriteSets[0]},
	"re1b_hig_def_b1_0": {L5_blush_open, spriteSets[0]},
	"re1b_kaii_b1_2": {smile_blush_close, spriteSets[0]},
	"re1b_komaru_b1_0": {fuan_blush_open, spriteSets[0]},
	"re1b_komaru_b2_0": {fuan_blush_open, spriteSets[0]},
	"re1b_okoru_b1_0": {sinken_blush_open, spriteSets[0]},
	"re1b_warai_b1_2": {smile_blush_close, spriteSets[0]},
	// ep2
	"re1b_nande_b1_1": {odoroki_blush_open, spriteSets[0]},
	// ep3
	"re1b_hig_okoru_b1_2": {sinken_open, spriteSets[0]},
	// ep5
	"re1b_def_b1_2": {smile_blush_open, spriteSets[0]},
	// ep6
	"re1b_bikkuri_b1_0": {odoroki_blush_open, spriteSets[0]},
	"re1b_hau_b1_0": {fuan_blush_open, spriteSets[0]},
	"re1b_kaii_b1_0": {smile_blush_close, spriteSets[0]},
	"re1b_warai_b1_0": {smile_blush_close, spriteSets[0]},
	// ep7
	"re1b_bikkuri_b1_1": {odoroki_blush_open, spriteSets[0]},
	"re1b_hau_b1_2": {fuan_blush_open, spriteSets[0]},
	"re1b_nande_b1_2": {odoroki_blush_open, spriteSets[0]},
	

	// rena 2a - normal casual
	// ep1
	"re2a_bikkuri_a1_2": {odoroki_blush_open, spriteSets[1]},
	"re2a_def_a1_0": {smile_blush_open, spriteSets[1]},
	"re2a_hau_a1_1": {fuan_blush_open, spriteSets[1]},
	"re2a_hig_def_a1_0": {L5_blush_open, spriteSets[1]},
	"re2a_hig_muhyou_a1_0": {L5_blush_open, spriteSets[1]},
	"re2a_kaii_a1_2": {smile_blush_close, spriteSets[1]},
	"re2a_komaru_a1_0": {fuan_blush_open, spriteSets[1]},
	"re2a_komaru_a2_0": {fuan_blush_open, spriteSets[1]},
	"re2a_nande_a1_1": {odoroki_blush_open, spriteSets[1]},
	"re2a_warai_a1_2": {smile_blush_close, spriteSets[1]},
	// ep2
	"re2a_okoru_a1_0": {sinken_blush_open, spriteSets[1]},
	// ep5
	"re2a_bikkuri_a1_1": {odoroki_blush_open, spriteSets[1]},
	"re2a_def_a1_2": {smile_blush_open, spriteSets[1]},
	"re2a_hau_a1_2": {fuan_blush_open, spriteSets[1]},
	"re2a_nande_a1_0": {odoroki_blush_open, spriteSets[1]},
	"re2a_warai_a1_1": {smile_blush_close, spriteSets[1]},
	// ep6
	"re2a_bikkuri_a1_0": {odoroki_blush_open, spriteSets[1]},
	"re2a_hau_a1_0": {fuan_blush_open, spriteSets[1]},
	"re2a_hig_okoru_a1_0": {sinken_open, spriteSets[1]},
	"re2a_kaii_a1_0": {smile_blush_close, spriteSets[1]},
	"re2a_warai_a1_0": {smile_blush_close, spriteSets[1]},
	// ep7
	"re2a_nande_a1_2": {odoroki_blush_open, spriteSets[1]},
	// rei
	"re2a_hig_okoru_a1_2": {sinken_open, spriteSets[1]},


	// rena 2b - hands casual
	// ep1
	"re2b_bikkuri_b1_2": {odoroki_blush_open, spriteSets[1]},
	"re2b_def_b1_0": {smile_blush_open, spriteSets[1]},
	"re2b_hau_b1_1": {fuan_blush_open, spriteSets[1]},
	"re2b_hig_def_b1_0": {L5_blush_open, spriteSets[1]},
	"re2b_hig_muhyou_b1_0": {L5_blush_open, spriteSets[1]},
	"re2b_hig_okoru_b1_0": {sinken_blush_open, spriteSets[1]},
	"re2b_kaii_b1_2": {smile_blush_close, spriteSets[1]},
	"re2b_komaru_b1_0": {fuan_blush_open, spriteSets[1]},
	"re2b_komaru_b2_1": {fuan_blush_open, spriteSets[1]},
	"re2b_warai_b1_2": {smile_blush_close, spriteSets[1]},
	// ep2
	"re2b_komaru_b2_0": {fuan_blush_open, spriteSets[1]},
	"re2b_nande_b1_1": {odoroki_blush_open, spriteSets[1]},
	"re2b_okoru_b1_0": {sinken_blush_open, spriteSets[1]},
	// ep3
	"re2b_hig_okoru_b1_2": {sinken_open, spriteSets[1]},
	// ep5
	"re2b_bikkuri_b1_1": {odoroki_blush_open, spriteSets[1]},
	"re2b_def_b1_2": {smile_blush_open, spriteSets[1]},
	"re2b_hau_b1_2": {fuan_blush_open, spriteSets[1]},
	"re2b_kaii_b1_0": {smile_blush_close, spriteSets[1]},
	"re2b_warai_b1_0": {smile_blush_close, spriteSets[1]},
	"re2b_warai_b1_1": {smile_blush_close, spriteSets[1]},
	// ep6
	"re2b_bikkuri_b1_0": {odoroki_blush_open, spriteSets[1]},
	"re2b_hau_b1_0": {fuan_blush_open, spriteSets[1]},
	"re2b_nande_b1_0": {odoroki_blush_open, spriteSets[1]},
	// ep7
	"re2b_nande_b1_2": {odoroki_blush_open, spriteSets[1]},

	// rena 3a - normal gym
	// ep1
	"re3a_bikkuri_a1_2": {odoroki_blush_open, spriteSets[38]},
	"re3a_def_a1_0": {smile_blush_open, spriteSets[38]},
	"re3a_hau_a1_1": {fuan_blush_open, spriteSets[38]},
	"re3a_kaii_a1_2": {smile_blush_close, spriteSets[38]},
	"re3a_komaru_a1_0": {fuan_blush_open, spriteSets[38]},
	"re3a_komaru_a2_0": {fuan_blush_open, spriteSets[38]},
	"re3a_nande_a1_1": {odoroki_blush_open, spriteSets[38]},
	"re3a_warai_a1_2": {smile_blush_close, spriteSets[38]},
	// ep6
	"re3a_hau_a1_0": {fuan_blush_open, spriteSets[38]},
	"re3a_kaii_a1_0": {smile_blush_close, spriteSets[38]},
	"re3a_okoru_a1_0": {sinken_blush_open, spriteSets[38]},
	"re3a_warai_a1_0": {smile_blush_close, spriteSets[38]},
	// rei
	"re3a_hau_a1_2": {fuan_blush_open, spriteSets[38]},

	// rena 3b - hands gym
	// ep1
	"re3b_bikkuri_b1_2": {odoroki_blush_open, spriteSets[38]},
	"re3b_hau_b1_1": {fuan_blush_open, spriteSets[38]},
	"re3b_kaii_b1_2": {smile_blush_close, spriteSets[38]},
	"re3b_warai_b1_2": {smile_blush_close, spriteSets[38]},
	// ep6
	"re3b_bikkuri_b1_0": {odoroki_blush_open, spriteSets[38]},
	"re3b_def_b1_0": {smile_blush_open, spriteSets[38]},
	"re3b_kaii_b1_0": {smile_blush_close, spriteSets[38]},
	"re3b_komaru_b1_0": {fuan_blush_open, spriteSets[38]},
	"re3b_nande_b1_0": {odoroki_blush_open, spriteSets[38]},
	"re3b_okoru_b1_0": {sinken_blush_open, spriteSets[38]},
	"re3b_warai_b1_0": {smile_blush_close, spriteSets[38]},

	// rena 6 - swimsuit
	// hou+
	"re6_bikkuri_a1_1": {odoroki_blush_open, spriteSets[31]},
	"re6_def_a1_2": {smile_blush_open, spriteSets[31]},
	"re6_hau_a1_2": {fuan_blush_open, spriteSets[31]},
	"re6_kaii_a1_2": {smile_blush_close, spriteSets[31]},
	"re6_komaru_a1_0": {fuan_blush_open, spriteSets[31]},
	"re6_nande_a1_2": {odoroki_blush_open, spriteSets[31]},
	"re6_warai_a1_2": {smile_blush_close, spriteSets[31]},

	// renasen 1 - nata attack
	// ep6
	"renasen1_def_0": {sinken_blush_open, spriteSets[3]},
	"renasen1_ikakaku_0": {sinken_blush_open, spriteSets[3]},
	"renasen1_muhyokaku_0": {L5_blush_open, spriteSets[3]},
	"renasen1_warai_0": {smile_blush_open, spriteSets[3]},
	// hou+
	"renasen1_tuukaku_0": {L5_blush_open, spriteSets[3]},


	// renasen 2 - nata down
	// ep6
	"renasen2_def_0": {sinken_blush_open, spriteSets[3]},
	"renasen2_ikakaku_0": {sinken_blush_open, spriteSets[3]},
	"renasen2_muhyokaku_0": {L5_blush_open, spriteSets[3]},
	"renasen2_shinken_0": {sinken_blush_open, spriteSets[3]},
	"renasen2_tuukaku_0": {L5_blush_open, spriteSets[3]},
	"renasen2_warai_0": {smile_blush_open, spriteSets[3]},

	// rika 1 - school
	// ep1
	"ri1_def_a1_0": {normal_blush_open, spriteSets[0]},
	"ri1_komaru_a1_0": {fuan_blush_open, spriteSets[0]},
	"ri1_niko_a1_0": {smile_blush_open, spriteSets[0]},
	"ri1_warai_a1_1": {smile_blush_close, spriteSets[0]},
	// ep2
	"ri1_fuman_a1_0": {sinken_blush_open, spriteSets[0]},
	"ri1_komaru_a2_0": {fuan_blush_open, spriteSets[0]},
	"ri1_majime_a1_1": {sinken_blush_open, spriteSets[0]},
	// ep3
	"ri1_majime_a1_0": {sinken_blush_open, spriteSets[0]},
	// ep5
	"ri1_warai_a1_2": {smile_blush_close, spriteSets[0]},
	// ep6
	"ri1_niyari_a1_0": {futeki_blush_open, spriteSets[0]},
	"ri1_warai_a1_0": {smile_blush_close, spriteSets[0]},
	// ep7
	"ri1_majime_a1_2": {sinken_blush_open, spriteSets[0]},
	"ri1_niko_a1_2": {smile_blush_open, spriteSets[0]},

	// rika 2 - casual
	// ep1
	"ri2_def_a1_0": {normal_blush_open, spriteSets[1]},
	"ri2_komaru_a1_0": {fuan_blush_open, spriteSets[1]},
	"ri2_niko_a1_0": {smile_blush_open, spriteSets[1]},
	// ep2
	"ri2_warai_a1_1": {smile_blush_close, spriteSets[1]},
	// ep3
	"ri2_komaru_a2_0": {fuan_blush_open, spriteSets[1]},
	// ep5
	"ri2_fuman_a1_0": {sinken_blush_open, spriteSets[1]},
	"ri2_niyari_a1_0": {futeki_blush_open, spriteSets[1]},
	"ri2_warai_a1_2": {smile_blush_close, spriteSets[1]},
	// ep6
	"ri2_warai_a1_0": {smile_blush_close, spriteSets[1]},
	// ep7
	"ri2_majime_a1_2": {sinken_blush_open, spriteSets[1]},
	"ri2_niko_a1_2": {smile_blush_open, spriteSets[1]},

	// rika 3 - gym
	// ep1
	"ri3_def_a1_0": {normal_blush_open, spriteSets[4]},
	"ri3_niko_a1_0": {smile_blush_open, spriteSets[4]},
	"ri3_warai_a1_1": {smile_blush_close, spriteSets[4]},
	// ep2
	"ri3_komaru_a1_0": {fuan_blush_open, spriteSets[4]},
	// ep6
	"ri3_warai_a1_0": {smile_blush_close, spriteSets[4]},

	// rika 4 - cat
	// ep1
	"ri4_komaru_a1_0": {fuan_blush_open, spriteSets[13]},
	"ri4_niko_a1_0": {smile_blush_open, spriteSets[13]},
	// rei
	"ri4_def_a1_0": {normal_blush_open, spriteSets[13]},
	"ri4_niko_a1_2": {smile_blush_open, spriteSets[13]},
	"ri4_warai_a1_2": {smile_blush_close, spriteSets[13]},

	// rika 5 - miko
	// ep1
	"ri5_def_a1_0": {normal_blush_open, spriteSets[10]},
	"ri5_komaru_a1_0": {fuan_blush_open, spriteSets[10]},
	"ri5_niko_a1_0": {smile_blush_open, spriteSets[10]},
	// ep2
	"ri5_warai_a1_1": {smile_blush_close, spriteSets[10]},
	// ep5
	"ri5_niko_a1_2": {smile_blush_open, spriteSets[10]},
	"ri5_warai_a1_2": {smile_blush_close, spriteSets[10]},

	// rika 6 - angel mort
	// ep6
	"ri6_def_a1_0": {normal_blush_open, spriteSets[5]},
	"ri6_komaru_a1_0": {fuan_blush_open, spriteSets[5]},
	"ri6_niko_a1_0": {smile_blush_open, spriteSets[5]},
	"ri6_warai_a1_0": {smile_blush_close, spriteSets[5]},
	// ep8
	"ri6_warai_a1_2": {smile_blush_close, spriteSets[5]},
	// rei
	"ri6_fuman_a1_0": {sinken_blush_open, spriteSets[5]},
	// hou+
	"ri6_komaru_a2_0": {fuan_blush_open, spriteSets[5]},

	// rika 8 - swimsuit
	// hou+
	"ri8_def_a1_0": {normal_blush_open, spriteSets[9]},
	"ri8_komaru_a1_0": {fuan_blush_open, spriteSets[9]},
	"ri8_komaru_a2_0": {fuan_blush_open, spriteSets[9]},
	"ri8_majime_a1_2": {fuan_open, spriteSets[9]},
	"ri8_niko_a1_2": {smile_blush_open, spriteSets[9]},
	"ri8_niyari_a1_0": {futeki_blush_open, spriteSets[9]},
	"ri8_warai_a1_2": {smile_blush_close, spriteSets[9]},

	// rika minor(?)
	// ep4
	"rim_def_0": {normal_blush_open, spriteSets[1]},
	"rim_komaru_0": {fuan_blush_open, spriteSets[1]},
	"rim_majime_0": {fuan_open, spriteSets[1]},
	"rim_niyari_0": {futeki_blush_open, spriteSets[1]},
	"rim_warai_0": {smile_blush_close, spriteSets[1]},
	"rim_warai_2": {smile_blush_close, spriteSets[1]},

	// satoko 1a - normal school
	// ep1
	"sa1a_akireru_a1_0": {normal_blush_open, spriteSets[0]},
	"sa1a_akuwarai_a1_1": {futeki_blush_open, spriteSets[0]},
	"sa1a_def_a1_1": {smile_blush_open, spriteSets[0]},
	"sa1a_hannbeso_a1_1": {odoroki_blush_open, spriteSets[0]},
	"sa1a_naku_a1_1": {odoroki_blush_close, spriteSets[0]},
	"sa1a_odoroki_a1_1": {sinken_blush_open, spriteSets[0]},
	"sa1a_warai_a1_1": {futeki_blush_close, spriteSets[0]},
	// ep2
	"sa1a_yareyare_a1_1": {normal_blush_close, spriteSets[0]},
	// ep3
	"sa1a_hannbeso_a3_0": {sinken_blush_open, spriteSets[0]},
	"sa1a_hau_a1_0": {smile_blush_open, spriteSets[0]},
	"sa1a_hau_a2_1": {smile_blush_open, spriteSets[0]},
	"sa1a_muhyou_a1_0": {L5_open, spriteSets[0]},
	"sa1a_muhyou_a2_0": {L5_open, spriteSets[0]},
	"sa1a_sakebu_a1_1": {odoroki_open, spriteSets[0]},
	"sa1a_warai_a1_0": {futeki_blush_close, spriteSets[0]},
	"sa1a_yareyare_a1_0": {normal_close, spriteSets[0]},
	"sa1a_yareyare_a2_0": {normal_blush_close, spriteSets[0]},
	// ep5
	"sa1a_hannbeso_a1_0": {odoroki_blush_open, spriteSets[0]},
	"sa1a_hannbeso_a3_2": {sinken_blush_open, spriteSets[0]},
	// ep6
	"sa1a_akuwarai_a1_0": {futeki_blush_open, spriteSets[0]},
	"sa1a_def_a1_0": {smile_blush_open, spriteSets[0]},
	"sa1a_odoroki_a1_0": {sinken_blush_open, spriteSets[0]},
	// ep7
	"sa1a_akuwarai_a1_2": {futeki_blush_open, spriteSets[0]},
	"sa1a_def_a1_2": {smile_blush_open, spriteSets[0]},
	"sa1a_hau_a2_2": {smile_blush_open, spriteSets[0]},
	"sa1a_muhyou_a2_2": {L5_open, spriteSets[0]},
	"sa1a_naku_a1_2": {odoroki_blush_close, spriteSets[0]},
	"sa1a_odoroki_a1_2": {sinken_blush_open, spriteSets[0]},
	"sa1a_sakebu_a1_2": {odoroki_open, spriteSets[0]},

	// satoko 1b - hands school
	// ep1
	"sa1b_akuwarai_b1_1": {futeki_blush_open, spriteSets[0]},
	"sa1b_def_b1_1": {smile_blush_open, spriteSets[0]},
	"sa1b_hannbeso_b1_1": {odoroki_blush_open, spriteSets[0]},
	"sa1b_hannbeso_b1_2": {odoroki_blush_open, spriteSets[0]},
	"sa1b_naku_b1_1": {odoroki_blush_close, spriteSets[0]},
	"sa1b_odoroki_b1_1": {sinken_blush_open, spriteSets[0]},
	"sa1b_odoroki_b1_2": {sinken_blush_open, spriteSets[0]},
	"sa1b_warai_b1_1": {futeki_blush_close, spriteSets[0]},
	"sa1b_yareyare_b2_1": {normal_blush_close, spriteSets[0]},
	// ep2
	"sa1b_akireru_b1_0": {normal_blush_open, spriteSets[0]},
	// ep3
	"sa1b_hannbeso_b3_0": {sinken_blush_open, spriteSets[0]},
	"sa1b_hau_b1_0": {smile_blush_open, spriteSets[0]},
	"sa1b_hau_b2_1": {smile_blush_open, spriteSets[0]},
	"sa1b_sakebu_b1_2": {odoroki_open, spriteSets[0]},
	"sa1b_yareyare_b1_0": {normal_close, spriteSets[0]},
	"sa1b_yareyare_b2_0": {normal_blush_close, spriteSets[0]},
	// ep5
	"sa1b_naku_b1_2": {normal_blush_open, spriteSets[0]},
	"sa1b_sakebu_b1_1": {odoroki_open, spriteSets[0]},
	"sa1b_warai_b1_0": {futeki_blush_close, spriteSets[0]},
	// ep6
	"sa1b_akuwarai_b1_0": {futeki_blush_open, spriteSets[0]},
	"sa1b_def_b1_0": {smile_blush_open, spriteSets[0]},
	"sa1b_muhyou_b1_0": {smile_blush_open, spriteSets[0]},
	"sa1b_muhyou_b2_0": {smile_blush_open, spriteSets[0]},
	"sa1b_odoroki_b1_0": {sinken_blush_open, spriteSets[0]},
	// ep7
	"sa1b_akuwarai_b1_2": {futeki_blush_open, spriteSets[0]},
	"sa1b_def_b1_2": {smile_blush_open, spriteSets[0]},
	"sa1b_hannbeso_b1_0": {sinken_blush_open, spriteSets[0]},
	"sa1b_muhyou_b2_2": {smile_blush_open, spriteSets[0]},

	// satoko 2a - normal casual
	// ep1
	"sa2a_akireru_a1_0": {normal_blush_open, spriteSets[1]},
	"sa2a_akuwarai_a1_1": {futeki_blush_open, spriteSets[1]},
	"sa2a_def_a1_1": {smile_blush_open, spriteSets[1]},
	"sa2a_hannbeso_a1_1": {odoroki_blush_open, spriteSets[1]},
	"sa2a_naku_a1_1": {odoroki_blush_close, spriteSets[1]},
	"sa2a_odoroki_a1_1": {sinken_blush_open, spriteSets[1]},
	"sa2a_warai_a1_1": {futeki_blush_close, spriteSets[1]},
	// ep2
	"sa2a_naku_a1_2": {odoroki_blush_close, spriteSets[1]},
	// ep3
	"sa2a_hau_a1_0": {smile_blush_open, spriteSets[1]},
	"sa2a_hau_a2_1": {smile_blush_open, spriteSets[1]},
	"sa2a_yareyare_a1_0": {normal_close, spriteSets[1]},
	"sa2a_yareyare_a2_0": {normal_blush_close, spriteSets[1]},
	// ep5
	"sa2a_hannbeso_a1_0": {odoroki_blush_open, spriteSets[1]},
	"sa2a_muhyou_a1_0": {smile_blush_open, spriteSets[1]},
	"sa2a_muhyou_a2_2": {smile_blush_open, spriteSets[1]},
	"sa2a_naku_a1_0": {odoroki_blush_close, spriteSets[1]},
	"sa2a_warai_a1_0": {futeki_blush_close, spriteSets[1]},
	"sa2a_def_a1_0": {smile_blush_open, spriteSets[1]},
	"sa2a_odoroki_a1_0": {sinken_blush_open, spriteSets[1]},
	// ep7
	"sa2a_akuwarai_a1_2": {futeki_blush_open, spriteSets[1]},
	"sa2a_def_a1_2": {smile_blush_open, spriteSets[1]},
	"sa2a_hannbeso_a3_2": {odoroki_blush_open, spriteSets[1]},
	"sa2a_hau_a1_2": {smile_blush_open, spriteSets[1]},
	"sa2a_odoroki_a1_2": {sinken_blush_open, spriteSets[1]},
	// ep8
	"sa2a_hau_a2_2": {smile_blush_open, spriteSets[1]},
	// hou+
	"sa2a_hannbeso_a1_2": {odoroki_blush_open, spriteSets[1]},

	// satoko 2b - hands casual
	// ep1
	"sa2b_akireru_b1_1": {normal_blush_open, spriteSets[1]},
	"sa2b_warai_b1_1": {futeki_blush_close, spriteSets[1]},
	"sa2b_yareyare_b1_1": {normal_blush_close, spriteSets[1]},
	// ep2
	"sa2b_akuwarai_b1_0": {futeki_blush_open, spriteSets[1]},
	"sa2b_def_b1_1": {smile_blush_open, spriteSets[1]},
	"sa2b_naku_b1_1": {odoroki_blush_close, spriteSets[1]},
	// ep3
	"sa2b_akireru_b1_0": {normal_blush_open, spriteSets[1]},
	"sa2b_akuwarai_b1_1": {futeki_blush_open, spriteSets[1]},
	"sa2b_hannbeso_b1_1": {odoroki_blush_open, spriteSets[1]},
	"sa2b_hau_b1_0": {smile_blush_open, spriteSets[1]},
	"sa2b_hau_b2_1": {smile_blush_open, spriteSets[1]},
	"sa2b_odoroki_b1_1": {sinken_blush_open, spriteSets[1]},
	"sa2b_yareyare_b1_0": {normal_close, spriteSets[1]},
	"sa2b_yareyare_b2_0": {normal_blush_close, spriteSets[1]},
	// ep5
	"sa2b_muhyou_b1_0": {smile_blush_open, spriteSets[1]},
	"sa2b_naku_b1_0": {odoroki_blush_close, spriteSets[1]},
	"sa2b_warai_b1_0": {futeki_blush_close, spriteSets[1]},
	// ep6
	"sa2b_odoroki_b1_0": {sinken_blush_open, spriteSets[1]},
	// ep7
	"sa2b_akuwarai_b1_2": {futeki_blush_open, spriteSets[1]},
	"sa2b_def_b1_2": {smile_blush_open, spriteSets[1]},
	"sa2b_hannbeso_b1_0": {odoroki_blush_open, spriteSets[1]},
	"sa2b_hau_b2_2": {smile_blush_open, spriteSets[1]},
	"sa2b_odoroki_b1_2": {sinken_blush_open, spriteSets[1]},
	// hou+
	"sa2b_sakebu_b1_2": {odoroki_open, spriteSets[1]},

	// satoko 3 - gym
	// ep1
	"sa3_akireru_a1_0": {normal_blush_open, spriteSets[48]},
	"sa3_akuwarai_a1_1": {futeki_blush_open, spriteSets[48]},
	"sa3_def_a1_1": {smile_blush_open, spriteSets[48]},
	"sa3_hannbeso_a1_1": {odoroki_blush_open, spriteSets[48]},
	"sa3_odoroki_a1_1": {sinken_blush_open, spriteSets[48]},
	"sa3_warai_a1_1": {futeki_blush_close, spriteSets[48]},
	// ep6
	"sa3_akuwarai_a1_0": {futeki_blush_open, spriteSets[48]},
	"sa3_def_a1_0": {smile_blush_open, spriteSets[48]},
	"sa3_hannbeso_a1_0": {odoroki_blush_open, spriteSets[48]},
	"sa3_odoroki_a1_0": {sinken_blush_open, spriteSets[48]},
	"sa3_warai_a1_0": {futeki_blush_close, spriteSets[48]},
	
	// satoko 4 - dog
	// ep1
	"sa4_akireru_a1_1": {normal_blush_open, spriteSets[0]},
	"sa4_odoroki_a1_1": {sinken_blush_open, spriteSets[0]},
	"sa4_warai_a1_1": {futeki_blush_close, spriteSets[0]},
	// rei
	"sa4_akireru_a1_0": {normal_blush_open, spriteSets[0]},

	// satoko 5 - towel
	// ep3
	"sa5_akireru_a1_0": {normal_open, spriteSets[7]},
	"sa5_hannbeso_a1_1": {odoroki_open, spriteSets[7]},
	"sa5_hannbeso_a3_1": {sinken_blush_open, spriteSets[7]},
	"sa5_hau_a1_0": {smile_open, spriteSets[7]},
	"sa5_odoroki_a1_1": {sinken_blush_open, spriteSets[7]},
	"sa5_sakebu_a1_1": {odoroki_open, spriteSets[7]},
	"sa5_warai_a1_1": {futeki_open, spriteSets[7]},
	"sa5_yareyare_a1_0": {normal_close, spriteSets[7]},
	"sa5_yareyare_a2_0": {normal_blush_close, spriteSets[7]},

	// satoko 6 - maid
	// ep6
	"sa6_akireru_a1_0": {normal_open, spriteSets[9]},
	"sa6_akuwarai_a1_0": {futeki_blush_open, spriteSets[9]},
	"sa6_hau_a1_0": {smile_open, spriteSets[9]},
	"sa6_odoroki_a1_0": {sinken_blush_open, spriteSets[9]},
	"sa6_warai_a1_0": {futeki_open, spriteSets[9]},
	"sa6_yareyare_a1_0": {normal_close, spriteSets[9]},
	// rei
	"sa6_yareyare_a2_0": {normal_blush_close, spriteSets[9]},

	// satoko 8a - blue dress
	// hou+
	"sa8a_akuwarai_a1_2": {futeki_blush_open, spriteSets[36]},
	"sa8a_def_a1_2": {smile_blush_open, spriteSets[36]},
	"sa8a_warai_a1_0": {futeki_open, spriteSets[36]},
	// satoko 9 - swimsuit
	// hou+
	"sa9_akireru_a1_0": {normal_open, spriteSets[2]},
	"sa9_hannbeso_a1_2": {odoroki_open, spriteSets[2]},
	"sa9_odoroki_a1_2": {sinken_blush_open, spriteSets[2]},
	"sa9_warai_a1_0": {futeki_open, spriteSets[2]},

	// satoko 10 - swimsuit 2
	// hou+
	"sa10_akireru_a1_0": {normal_open, spriteSets[2]},
	"sa10_akuwarai_a1_2": {futeki_blush_open, spriteSets[2]},
	"sa10_def_a1_2": {smile_blush_open, spriteSets[2]},
	"sa10_muhyou_a2_2": {smile_blush_open, spriteSets[2]},
	"sa10_odoroki_a1_2": {sinken_blush_open, spriteSets[2]},
	"sa10_warai_a1_0": {futeki_open, spriteSets[2]},
	"sa10_yareyare_a1_0": {normal_close, spriteSets[2]},
	"sa10_yareyare_a2_0": {normal_blush_close, spriteSets[2]},

	// satoko 11 - towel 2
	// hou+
	"sa11_akireru_a1_0": {normal_open, spriteSets[7]},
	"sa11_odoroki_a1_2": {sinken_blush_open, spriteSets[7]},
	"sa11_warai_a1_0": {futeki_open, spriteSets[7]},
	"sa11_yareyare_a1_0": {normal_close, spriteSets[7]},

	// takano 1 - casual
	// ep1
	"ta1_akuwarai_1": {futeki_open, spriteSets[0]},
	"ta1_def_0": {smile_open, spriteSets[0]},
	"ta1_def_1": {smile_open, spriteSets[0]},
	"ta1_warai_1": {smile_close, spriteSets[0]},
	// ep2
	"ta1_hatena_0": {smile_open, spriteSets[0]},
	"ta1_human_0": {futeki_open, spriteSets[0]},
	// ep3
	"ta1_hatena_1": {smile_open, spriteSets[0]},
	"ta1_human_1": {futeki_open, spriteSets[0]},
	// ep5
	"ta1_warai_2": {smile_close, spriteSets[0]},
	// ep6
	"ta1_akuwarai_0": {futeki_open, spriteSets[0]},
	"ta1_warai_0": {smile_close, spriteSets[0]},
	// ep7
	"ta1_akuwarai_2": {futeki_open, spriteSets[0]},
	// ep8
	"ta1_iradachi_0": {sinken_open, spriteSets[0]},
	"ta1_kanashimi_0": {fuan_open, spriteSets[0]},
	"ta1_sakebi_2": {sinken_open, spriteSets[0]},

	// takano 2 - nurse
	// ep7
	"ta2_akuwarai_2": {futeki_open, spriteSets[1]},
	"ta2_def_0": {smile_open, spriteSets[1]},
	"ta2_hatena_0": {smile_open, spriteSets[1]},
	"ta2_human_0": {futeki_open, spriteSets[1]},
	"ta2_warai_2": {smile_close, spriteSets[1]},
	// ep8
	"ta2_iradachi_2": {sinken_open, spriteSets[1]},
	"ta2_kanashimi_0": {fuan_open, spriteSets[1]},
	"ta2_sakebi_0": {sinken_open, spriteSets[1]},
	"ta2_sakebi_2": {sinken_open, spriteSets[1]},

	// takano 3 - army
	// ep7
	"ta3_akuwarai_2": {futeki_open, spriteSets[9]},
	"ta3_def_0": {smile_open, spriteSets[9]},
	// ep8
	"ta3_human_0": {futeki_open, spriteSets[9]},
	"ta3_iradachi_0": {sinken_open, spriteSets[9]},
	"ta3_sakebi_2": {sinken_open, spriteSets[9]},
	// hou+
	"ta3_hatena_0": {smile_open, spriteSets[9]},
	"ta3_warai_2": {smile_close, spriteSets[9]},

	// takano 5 - army hatless
	// ep8
	"ta5_akuwarai_2": {futeki_open, spriteSets[10]},
	"ta5_human_0": {futeki_open, spriteSets[10]},
	"ta5_iradachi_0": {sinken_open, spriteSets[10]},
	"ta5_sakebi_2": {sinken_open, spriteSets[10]},

	// takano 7 - army bunny
	// ep8
	"ta7_hatena_0": {smile_open, spriteSets[10]},
	"ta7_sakebi_2": {sinken_open, spriteSets[10]},


	// chie 1
	// ep1
	"tie_1_0": {smile_open, spriteSets[0]},
	"tie_2_0": {sinken_open, spriteSets[0]},
	// ep2
	"tie_3_1": {fuan_open, spriteSets[0]},
	"tie_4_0": {sinken_open, spriteSets[0]},
	// ep6
	"tie_3_0": {fuan_open, spriteSets[0]},
	// ep7
	"tie_3_2": {fuan_open, spriteSets[0]},

	// tomitake 1 - casual
	// ep1
	"tomi1_def_0": {smile_open, spriteSets[0]},
	"tomi1_komaru_1": {fuan_open, spriteSets[0]},
	"tomi1_warai_1": {smile_close, spriteSets[0]},
	// ep5
	"tomi1_warai_2": {smile_close, spriteSets[0]},
	// ep6
	"tomi1_komaru_0": {fuan_open, spriteSets[0]},
	"tomi1_warai_0": {smile_close, spriteSets[0]},
	"tomi3_def_0": {smile_open, spriteSets[0]},
	// ep7
	"tomi1_komaru_2": {fuan_open, spriteSets[0]},
	// ep8
	"tomi1_shinken_0": {sinken_open, spriteSets[0]},
	"tomi1_shinken_2": {sinken_open, spriteSets[0]},
	// rei
	"tomi1_ikari_2": {sinken_open, spriteSets[0]},

	// tomitake 2 - army
	// ep7
	"tomi2_def_0": {smile_open, spriteSets[4]},
	"tomi2_komaru_2": {fuan_open, spriteSets[4]},
	"tomi2_warai_2": {smile_close, spriteSets[4]},
	// ep8
	"tomi2_shinken_0": {sinken_open, spriteSets[4]},


	// tomitake 3 - casual?
	"tomi3_ikari_2": {sinken_open, spriteSets[0]},
	"tomi3_komaru_2": {sinken_open, spriteSets[0]},
	"tomi3_shinken_2": {sinken_open, spriteSets[0]},
	"tomi3_warai_2": {smile_close, spriteSets[0]},
	// hou+
	"tomi3_shinken_0": {sinken_open, spriteSets[0]},
	

	// kasai
	// ep2
	"kasa_1_0": {smile_open, spriteSets[0]},
	"kasa_2_0": {odoroki_open, spriteSets[0]},
	// ep5
	"kasa_2_2": {odoroki_open, spriteSets[0]},
	"kasa_3_0": {sinken_open, spriteSets[0]},

	// shion 1a - normal casual
	// ep2
	"si1a_akuwarai_a1_2": {futeki_blush_open, spriteSets[1]},
	"si1a_def_a1_0": {smile_blush_open, spriteSets[1]},
	"si1a_hau_a1_1": {fuan_blush_open, spriteSets[1]},
	"si1a_huteki_a1_1": {futeki_blush_open, spriteSets[1]},
	"si1a_ikari_a1_2": {sinken_blush_open, spriteSets[1]},
	"si1a_majime_a1_0": {sinken_blush_open, spriteSets[1]},
	"si1a_odoroki_a1_2": {odoroki_blush_open, spriteSets[1]},
	"si1a_warai_a1_2": {smile_blush_close, spriteSets[1]},
	"si1a_wink_a1_2": {smile_blush_close, spriteSets[1]},
	"si1a_yowaki_a1_1": {fuan_blush_open, spriteSets[1]},
	// ep3
	"si1a_tohoho_a1_0": {normal_blush_open, spriteSets[1]},
	"si1a_tokui_a1_2": {futeki_blush_close, spriteSets[1]},
	// ep6
	"si1a_akuwarai_a1_0": {futeki_blush_open, spriteSets[1]},
	"si1a_odoroki_a1_0": {odoroki_blush_open, spriteSets[1]},
	"si1a_tokui_a1_0": {futeki_blush_close, spriteSets[1]},
	"si1a_warai_a1_0": {smile_blush_close, spriteSets[1]},
	"si1a_wink_a1_0": {smile_blush_close, spriteSets[1]},
	// ep7
	"si1a_tokui_a1_1": {futeki_blush_close, spriteSets[1]},
	"si1a_yowaki_a1_2": {fuan_blush_open, spriteSets[1]},
	// hou+
	"si1a_huteki_a1_2": {futeki_blush_open, spriteSets[1]},

	// shion 1b - hand casual
	// ep2
	"si1b_akuwarai_b1_2": {futeki_blush_open, spriteSets[1]},
	"si1b_def_b1_0": {smile_blush_open, spriteSets[1]},
	"si1b_hau_b1_1": {fuan_blush_open, spriteSets[1]},
	"si1b_huteki_b1_1": {futeki_blush_open, spriteSets[1]},
	"si1b_tokui_b1_2": {futeki_blush_close, spriteSets[1]},
	"si1b_warai_b1_2": {smile_blush_close, spriteSets[1]},
	"si1b_wink_b1_2": {smile_blush_close, spriteSets[1]},
	// ep3
	"si1b_majime_b1_0": {sinken_blush_open, spriteSets[1]},
	// ep6
	"si1b_wink_b1_0": {smile_blush_close, spriteSets[1]},
	// ep8
	"si1b_odoroki_b1_2": {odoroki_blush_open, spriteSets[1]},
	"si1b_tohoho_b1_0": {normal_blush_open, spriteSets[1]},
	// rei
	"si1b_yowaki_b1_2": {fuan_blush_open, spriteSets[1]},
	// hou+
	"si1b_tokui_b1_1": {futeki_blush_close, spriteSets[1]},


	// shion 2 - angel mort
	// ep2
	"si2_akuwarai_a1_2": {futeki_blush_open, spriteSets[3]},
	"si2_def_a1_0": {smile_blush_open, spriteSets[3]},
	"si2_hau_a1_1": {fuan_blush_open, spriteSets[3]},
	"si2_huteki_a1_2": {futeki_blush_open, spriteSets[3]},
	"si2_majime_a1_0": {sinken_blush_open, spriteSets[3]},
	"si2_odoroki_a1_2": {odoroki_blush_open, spriteSets[3]},
	"si2_tokui_a1_2": {futeki_blush_close, spriteSets[3]},
	"si2_warai_a1_2": {smile_blush_close, spriteSets[3]},
	"si2_wink_a1_2": {smile_blush_close, spriteSets[3]},
	"si2_yowaki_a1_1": {fuan_blush_open, spriteSets[3]},
	// ep6
	"si2_akuwarai_a1_0": {futeki_blush_open, spriteSets[3]},
	"si2_warai_a1_0": {smile_blush_close, spriteSets[3]},
	"si2_wink_a1_0": {smile_blush_close, spriteSets[3]},
	// hou+
	"si2_tohoho_a1_0": {normal_blush_open, spriteSets[3]},
	"si2_yowaki_a1_2": {fuan_blush_open, spriteSets[3]},


	// shion 3 - school
	// ep5
	"si3_akuwarai_a1_1": {futeki_blush_open, spriteSets[0]},
	"si3_def_a1_0": {smile_blush_open, spriteSets[0]},
	"si3_tokui_a1_1": {futeki_blush_close, spriteSets[0]},
	"si3_warai_a1_1": {smile_blush_close, spriteSets[0]},
	"si3_wink_a1_0": {smile_blush_close, spriteSets[0]},
	// ep7
	"si3_akuwarai_a1_2": {futeki_blush_open, spriteSets[0]},
	"si3_huteki_a1_2": {futeki_blush_open, spriteSets[0]},
	"si3_ikari_a1_2": {sinken_blush_open, spriteSets[0]},
	"si3_majime_a1_0": {sinken_blush_open, spriteSets[0]},
	"si3_odoroki_a1_2c": {odoroki_blush_open, spriteSets[0]},
	"si3_tohoho_a1_0": {normal_blush_open, spriteSets[0]},
	"si3_warai_a1_2": {smile_blush_close, spriteSets[0]},
	"si3_yowaki_a1_2": {fuan_blush_open, spriteSets[0]},
	// rei
	"si3_hau_a1_1": {fuan_blush_open, spriteSets[0]},
	"si3_huteki_a1_1": {futeki_blush_open, spriteSets[0]},
	"si3_tokui_a1_2": {futeki_blush_close, spriteSets[0]},

	// shion 5 - damaged work
	// hou+	
	"si5_akuwarai_a1_2": {futeki_blush_open, spriteSets[7]},
	"si5_huteki_a1_2": {futeki_blush_open, spriteSets[7]},
	"si5_majime_a1_0": {sinken_blush_open, spriteSets[7]},
	"si5_odoroki_a1_2": {odoroki_blush_open, spriteSets[7]},
	"si5_tokui_a1_1": {futeki_blush_close, spriteSets[7]},
	// shion 6 - maid work
	// hou+	
	"si6_akuwarai_a1_2": {futeki_blush_open, spriteSets[2]},
	"si6_def_a1_0": {smile_blush_open, spriteSets[2]},
	"si6_huteki_a1_2": {futeki_blush_open, spriteSets[2]},
	"si6_ikari_a1_2": {sinken_blush_open, spriteSets[2]},
	"si6_majime_a1_0": {sinken_blush_open, spriteSets[2]},
	"si6_odoroki_a1_2": {odoroki_blush_open, spriteSets[2]},
	"si6_tohoho_a1_0": {normal_blush_open, spriteSets[2]},
	"si6_tokui_a1_1": {futeki_blush_close, spriteSets[2]},
	"si6_warai_a1_2": {smile_blush_close, spriteSets[2]},
	"si6_wink_a1_0": {smile_blush_close, spriteSets[2]},

	// irie 1 - casual
	// ep3
	"iri1_def1_0": {futeki_open, spriteSets[3]},
	"iri1_def2_1": {smile_open, spriteSets[3]},
	"iri1_majime_1": {normal_open, spriteSets[3]},
	"iri1_majime2_1": {sinken_open, spriteSets[3]},
	"iri1_majime3_1": {normal_open, spriteSets[3]},
	"iri1_warai_2": {smile_open, spriteSets[3]},
	// ep5
	"iri1_def2_2": {smile_open, spriteSets[3]},
	"iri1_majime_2": {normal_open, spriteSets[3]},
	"iri1_majime2_0": {sinken_open, spriteSets[3]},
	// ep6
	"iri1_majime_0": {normal_open, spriteSets[3]},
	"iri1_warai_0": {smile_open, spriteSets[3]},

	// irie 2 - doctor
	// ep3
	"iri2_def1_0": {futeki_open, spriteSets[0]},
	"iri2_def2_1": {smile_open, spriteSets[0]},
	"iri2_majime_1": {normal_open, spriteSets[0]},
	"iri2_majime2_1": {sinken_open, spriteSets[0]},
	"iri2_majime3_1": {normal_open, spriteSets[0]},
	"iri2_warai_2": {smile_open, spriteSets[0]},
	// ep5
	"iri2_majime_2": {normal_open, spriteSets[0]},
	// ep6
	"iri2_def2_0": {smile_open, spriteSets[0]},
	"iri2_majime_0": {normal_open, spriteSets[0]},
	"iri2_majime2_0": {sinken_open, spriteSets[0]},
	"iri2_warai_0": {smile_open, spriteSets[0]},
	// ep7
	"iri2_def2_2": {smile_open, spriteSets[0]},

	// irie 3 - coach
	// ep3	
	"iri3_def1_0": {futeki_open, spriteSets[2]},
	"iri3_def2_1": {smile_open, spriteSets[2]},
	"iri3_warai_2": {smile_open, spriteSets[2]},
	// ep5
	"iri3_def2_2": {smile_open, spriteSets[2]},
	"iri3_majime_2": {normal_open, spriteSets[2]},
	"iri3_majime2_0": {sinken_open, spriteSets[2]},
	// ep6
	"iri3_def2_0": {smile_open, spriteSets[2]},
	"iri3_majime_0": {normal_open, spriteSets[2]},
	"iri3_warai_0": {smile_open, spriteSets[2]},

	// chibi mion
	// ch4
	"chibimion_def_0": {smile_open, spriteSets[1]},
	"chibimion_def_2": {smile_open, spriteSets[1]},
	"chibimion_shinken_0": {normal_open, spriteSets[1]},
	"chibimion_warai_1": {smile_close, spriteSets[1]},
	"chibimion_warai_2": {smile_close, spriteSets[1]},
	"chibimion_wink_0": {smile_close, spriteSets[1]},
	"chibimion_wink_1": {smile_close, spriteSets[1]},

	// akane
	// ep5
	"aka_def_0": {normal_open, spriteSets[0]},
	"aka_sakebi_0": {normal_open, spriteSets[0]},
	"aka_warai_0": {normal_open, spriteSets[0]},

	// k1 1 - school
	// ep5
	"kei1_def1_0": {smile_open, spriteSets[0]},
	"kei1_def2_0": {futeki_open, spriteSets[0]},
	"kei1_ikari1_0": {sinken_open, spriteSets[0]},
	"kei1_komaru_0": {fuan_open, spriteSets[0]},
	"kei1_majime_0": {normal_open, spriteSets[0]},
	"kei1_majime2_0": {normal_open, spriteSets[0]},
	"kei1_nayamu_2": {sinken_close, spriteSets[0]},
	"kei1_warai_2": {smile_close, spriteSets[0]},
	// ep6
	"kei1_warai_0": {smile_close, spriteSets[0]},
	// ep7
	"kei1_ikari2_2": {sinken_blush_open, spriteSets[0]},

	// k1 2 - casual
	// ep5
	"kei2_def1_0": {smile_open, spriteSets[1]},
	"kei2_def2_0": {futeki_open, spriteSets[1]},
	"kei2_ikari1_0": {sinken_open, spriteSets[1]},
	"kei2_ikari2_1": {sinken_blush_open, spriteSets[1]},
	"kei2_komaru_0": {fuan_open, spriteSets[1]},
	"kei2_majime_0": {normal_open, spriteSets[1]},
	"kei2_majime2_0": {normal_open, spriteSets[1]},
	"kei2_warai_2": {smile_close, spriteSets[1]},
	// ep6
	"kei2_nayamu_0": {sinken_close, spriteSets[1]},
	"kei2_warai_0": {smile_close, spriteSets[1]},
	// ep7
	"kei2_ikari2_2": {sinken_open, spriteSets[1]},
	"kei2_nayamu_2": {sinken_close, spriteSets[1]},
	// hou+
	"kei2_hig_0": {L5_open, spriteSets[1]},
	"kei2_hig2_2": {L5_open, spriteSets[1]},

	// k1 5 - casual bat
	// hou+
	"kei5_def1_0": {smile_open, spriteSets[1]},
	"kei5_def2_0": {futeki_open, spriteSets[1]},
	"kei5_hig_0": {L5_open, spriteSets[1]},
	"kei5_ikari1_0": {sinken_open, spriteSets[1]},
	"kei5_ikari2_2": {sinken_blush_open, spriteSets[1]},
	"kei5_komaru_0": {fuan_open, spriteSets[1]},
	"kei5_nayamu_2": {sinken_close, spriteSets[1]},
	"kei5_warai_2": {smile_close, spriteSets[1]},


	// k1 6 - maid
	// rei
	"kei6_komaru_0": {fuan_open, spriteSets[27]},
	"kei6_nayamu_2": {sinken_close, spriteSets[27]},

	// k1 7 - swimsuit
	// hou+
	"kei7_def1_0": {smile_open, spriteSets[18]},
	"kei7_def2_0": {futeki_open, spriteSets[18]},
	"kei7_ikari1_0": {sinken_open, spriteSets[18]},
	"kei7_ikari2_2": {sinken_blush_open, spriteSets[18]},
	"kei7_komaru_0": {fuan_open, spriteSets[18]},
	"kei7_majime_0": {normal_open, spriteSets[18]},
	"kei7_majime2_0": {normal_open, spriteSets[18]},
	"kei7_nayamu_2": {sinken_close, spriteSets[18]},
	"kei7_warai_2": {smile_close, spriteSets[18]},

	// k1 8 - girl
	// hou+
	"kei8_ikari1_0": {sinken_open, spriteSets[13]},
	"kei8_ikari2_2": {sinken_blush_open, spriteSets[13]},
	"kei8_komaru_0": {fuan_open, spriteSets[13]},
	"kei8_majime_0": {normal_open, spriteSets[13]},
	"kei8_nayamu_2": {sinken_close, spriteSets[13]},
	"kei8_warai_2": {smile_close, spriteSets[13]},

	// keisen - casual bat?
	// hou+
	"keisen_niramu_0": {sinken_blush_open, spriteSets[1]},
	"keisen_shinken_0": {normal_open, spriteSets[1]},

	// satoshi 1 - casual
	// ep5
	"sato1_def1_0": {smile_open, spriteSets[0]},
	"sato1_def2_0": {smile_close, spriteSets[0]},
	"sato1_ikari_1": {sinken_open, spriteSets[0]},
	"sato1_komaru_0": {fuan_open, spriteSets[0]},
	"sato1_komaru2_0": {fuan_open, spriteSets[0]},
	"sato1_tukare_0": {fuan_close, spriteSets[0]},
	"sato1_warai_0": {smile_open, spriteSets[0]},
	"sato1_warai_1": {smile_open, spriteSets[0]},

	// satoshi 2 - baseball
	// ep5
	"sato2_def1_0": {smile_open, spriteSets[1]},
	"sato2_def2_0": {smile_blush_close, spriteSets[1]},
	"sato2_komaru_0": {fuan_open, spriteSets[1]},
	"sato2_komaru2_0": {fuan_open, spriteSets[1]},
	"sato2_tukare_0": {fuan_close, spriteSets[1]},
	"sato2_warai_1": {smile_open, spriteSets[1]},

	// teppei
	// ep5
	"tetu_1_0": {futeki_open, spriteSets[0]},
	// ep6
	"tetu_2_0": {normal_open, spriteSets[0]},
	"tetu_3_0": {odoroki_open, spriteSets[0]},
	// ep7
	"tetu_2_2": {normal_open, spriteSets[0]},
	"tetu_3_2": {odoroki_open, spriteSets[0]},
	// hou+
	"tetu_4_2": {smile_open, spriteSets[0]},
	"tetu_5_0": {smile_open, spriteSets[0]},
	
	// rina 1
	// rina never appeared in mei
	// ep6
	"rina_def_0": {smile_open, spriteSets[0]},
	"rina_ikari_0": {sinken_open, spriteSets[0]},
	"rina_warai_0": {smile_open, spriteSets[0]},
	// ep7
	"rina_warai_2": {smile_open, spriteSets[0]},

	// akasaka 1 - casual
	// ep7
	"aks1_def_0": {smile_open, spriteSets[0]},
	"aks1_warai_2": {smile_close, spriteSets[0]},
	// ep8
	"aks1_sakebi_2": {sinken_open, spriteSets[0]},
	"aks1_shinken_0": {normal_open, spriteSets[0]},

	// akasaka 2 - fighting
	// ep8
	"aks2_niyari_0": {smile_open, spriteSets[1]},
	"aks2_sakebi_2": {sinken_open, spriteSets[1]},
	"aks2_shinken_0": {normal_open, spriteSets[1]},


	// hanyuu 1 - formal
	// ep7
	"ha1_au_2": {fuan_blush_open, spriteSets[8]},
	"ha1_def_0": {smile_blush_open, spriteSets[8]},
	"ha1_def2_0": {normal_blush_open, spriteSets[8]},
	"ha1_odoroki_2": {odoroki_blush_open, spriteSets[8]},
	"ha1_warai_2": {smile_blush_close, spriteSets[8]},
	"ha1_yowaki_0": {normal_blush_open, spriteSets[8]},
	// ep8
	"ha1_muhyou_0": {normal_blush_open, spriteSets[8]},
	"ha1_sakebi_0": {sinken_blush_open, spriteSets[8]},
	"ha1_shinken_0": {sinken_open, spriteSets[8]},

	// hanyuu 2a - school
	// ep8
	"ha2a_au_2": {fuan_blush_open, spriteSets[0]},
	"ha2a_def_0": {smile_blush_open, spriteSets[0]},
	"ha2a_def2_0": {sinken_blush_open, spriteSets[0]},
	"ha2a_muhyou_0": {normal_blush_open, spriteSets[0]},
	"ha2a_odoroki_2": {odoroki_blush_open, spriteSets[0]},
	"ha2a_sakebi_0": {sinken_blush_open, spriteSets[0]},
	"ha2a_warai_2": {smile_blush_close, spriteSets[0]},
	"ha2a_yowaki_0": {normal_blush_open, spriteSets[0]},

	// hanyuu 2b - school
	// ep8
	"ha2b_def_0": {smile_blush_open, spriteSets[0]},
	"ha2b_def2_0": {normal_blush_open, spriteSets[0]},
	"ha2b_warai_2": {smile_blush_close, spriteSets[0]},

	// hanyuu 3a - school
	// ep8
	"ha3a_au_2": {fuan_blush_open, spriteSets[0]},
	"ha3a_def_0": {smile_blush_open, spriteSets[0]},
	"ha3a_def2_0": {normal_blush_open, spriteSets[0]},
	"ha3a_odoroki_2": {odoroki_blush_open, spriteSets[0]},
	"ha3a_warai_2": {smile_blush_close, spriteSets[0]},
	"ha3a_yowaki_0": {normal_blush_open, spriteSets[0]},
	// rei
	"ha3a_shinken_0": {sinken_blush_open, spriteSets[0]},

	// hanyuu 5 - costume
	// hou+
	"ha5_muhyou_0": {normal_blush_open, spriteSets[31]},
	"ha5_odoroki_2": {odoroki_blush_open, spriteSets[31]},
	"ha5_shinken_0": {sinken_blush_open, spriteSets[31]},

	// hanyuu 6 - angel mort
	// hou+
	"ha6_au_2": {fuan_blush_open, spriteSets[5]},

	// okonogi 1 - casual
	// okonogi never appeared in mei
	// ep7
	"oko_def_0": {smile_open, spriteSets[0]},
	// ep8
	"oko_kumon_0": {futeki_open, spriteSets[0]},
	"oko_niyari_2": {sinken_open, spriteSets[0]},
	"oko_odoroki_0": {normal_open, spriteSets[0]},
	"oko_sakebi_0": {fuan_open, spriteSets[0]},

	// okonogi 2 - army
	// ep8
	"oko2_def_0": {smile_open, spriteSets[1]},

	// okonogi 3 - army?
	// rei
	"oko3_def_0": {smile_open, spriteSets[2]},
	"oko3_kumon_2": {futeki_open, spriteSets[2]},
	"oko3_niyari_2": {sinken_open, spriteSets[2]},
	"oko3_odoroki_0": {normal_open, spriteSets[2]},
	"oko3_sakebi_1": {fuan_open, spriteSets[2]},

	// kameda 1a - speedo
	// kameda never appeared in mei
	// hou+
	"kameda1a_def_0": {smile_open, spriteSets[0]},
	"kameda1a_shinken_0": {normal_open, spriteSets[0]},
	"kameda1a_warai_2": {smile_close, spriteSets[0]},
	"kameda1b_odoroki_2": {odoroki_open, spriteSets[0]},

	// mo 1-3 & mura
	// i have no idea who these characters are
	// hou+
	"mo1_01_0": {normal_open, spriteSets[0]},
	"mo2_01_0": {smile_open, spriteSets[0]},
	"mo3_01_0": {fuan_open, spriteSets[0]},
	"mura_01_0": {normal_open, spriteSets[0]},

	// tamura 1
	// hou+
	"tamura1a_01_0": {normal_blush_open, spriteSets[0]},
	"tamura1a_02_2": {fuan_blush_open, spriteSets[0]},
	"tamura1a_03_2": {fuan_blush_open, spriteSets[0]},
	"tamura1a_04_2": {fuan_blush_close, spriteSets[0]},
	"tamura1a_05_0": {normal_blush_close, spriteSets[0]},
	"tamura1a_06_0": {futeki_blush_close, spriteSets[0]},
	"tamura1a_07_2": {odoroki_blush_open, spriteSets[0]},
	"tamura1a_08_0": {sinken_blush_open, spriteSets[0]},
	"tamura1a_09_2": {odoroki_blush_open, spriteSets[0]},
	"tamura1a_10_2": {odoroki_blush_close, spriteSets[0]},
	"tamura1a_11_2": {fuan_open, spriteSets[0]},

	// tamura 2
	// hou+
	"tamura2a_10_2": {odoroki_blush_close, spriteSets[2]},

	// une 1a
	// hou+
	"une1a_01_0": {normal_open, spriteSets[0]},
	"une1a_02_2": {odoroki_open, spriteSets[0]},
	"une1a_03_0": {normal_close, spriteSets[0]},
	"une1a_04_2": {smile_open, spriteSets[0]},
	"une1a_05_0": {fuan_close, spriteSets[0]},
	"une1a_06_0": {normal_open, spriteSets[0]},
	"une1a_07_1": {fuan_blush_open, spriteSets[0]},
	"une1a_09_2": {odoroki_blush_open, spriteSets[0]},
	"une1a_10_2": {L5_blush_close, spriteSets[0]},
	"une1a_11_1": {odoroki_open, spriteSets[0]},
	"une1a_12_0": {futeki_blush_open, spriteSets[0]},
	"une1a_14_2": {futeki_blush_open, spriteSets[0]},
	"une1a_15_0": {L5_blush_open, spriteSets[0]},
	"une1b_01_0": {normal_open, spriteSets[0]},
	"une1b_04_2": {smile_open, spriteSets[0]},
	"une1b_06_0": {normal_open, spriteSets[0]},
	"une1b_07_1": {fuan_blush_open, spriteSets[0]},
	"une1b_09_2": {odoroki_blush_open, spriteSets[0]},
	"une1b_10_2": {L5_blush_close, spriteSets[0]},
	"une1b_11_1": {odoroki_open, spriteSets[0]},
	"une1b_12_0": {futeki_blush_open, spriteSets[0]},
	"une1b_13_2": {futeki_blush_open, spriteSets[0]},
	"une2b_10_2": {L5_blush_close, spriteSets[0]},
	"une3a_01_0": {normal_open, spriteSets[0]},
	"une3a_02_2": {odoroki_open, spriteSets[0]},
	"une3a_05_0": {fuan_close, spriteSets[0]},
	"une3a_06_0": {normal_open, spriteSets[0]},
	"une3a_07_1": {fuan_blush_open, spriteSets[0]},
	"une3a_08_2": {odoroki_open, spriteSets[0]},
	"une3a_09_2": {odoroki_blush_open, spriteSets[0]},
	"une3a_10_2": {L5_blush_close, spriteSets[0]},
	"une3a_11_1": {odoroki_open, spriteSets[0]},
	"une3a_12_0": {futeki_blush_open, spriteSets[0]},
	"une3a_13_2": {futeki_blush_open, spriteSets[0]},
	"une3a_14_2": {futeki_blush_open, spriteSets[0]},
	"une3a_15_0": {L5_blush_open, spriteSets[0]},
	"une4a_01_0": {normal_open, spriteSets[0]},
	"une4a_02_2": {odoroki_open, spriteSets[0]},
	"une4a_09_2": {odoroki_blush_open, spriteSets[0]},
}


// Converts a sprite key into a folder name using FolderMap
func GetFolder(key string) string {
	var selected string
	longest := 0
	for prefix, folder := range FolderMap {
		if strings.HasPrefix(key, prefix) && len(prefix) > longest {
			selected = folder
			longest = len(prefix)
		}
	}
	if selected == "" {
		log.Printf("WARNING: No folder mapping found for key: %s", key)
		selected = "unknown"
	}
	return selected
}

// Returns the first existing sprite file path while checking fallback variants
func ResolveSpritePath(key string) string {
	info, ok := RawGameSprites[key]
	if !ok {
		log.Fatalf("Sprite key not found: %s", key)
	}
	expression := info[0]
	preferredVariant := info[1]
	folder := GetFolder(key)

	// Build fallback variant list (descending v006 → v001)
	tryVariants := []string{preferredVariant, "v006", "v005", "v004", "v003", "v002", "v001"}

	seen := make(map[string]bool)
	for _, v := range tryVariants {
		if seen[v] {
			continue
		}
		seen[v] = true

		candidate := filepath.Join("sprites", "mei", folder, v, expression+".png")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Fallback if nothing exists
	fallback := filepath.Join("sprites", "mei", folder, preferredVariant, expression+".png")
	log.Printf("WARNING: No variant found for %s/%s — using fallback: %s", folder, expression, fallback)
	return fallback
}