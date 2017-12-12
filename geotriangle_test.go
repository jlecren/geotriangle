package geotriangle

import "testing"
import "github.com/stretchr/testify/assert"

func TestBuildCodeMask( t *testing.T ) {

	assert.Equal(t, uint64(0xFC00000000000000), buildPathMask(0), "A depth 0 the mask is 0xFC00000000000000")
	assert.Equal(t, uint64(0xFF00000000000000), buildPathMask(1), "A depth 1 the mask is 0xFF00000000000000")
	assert.Equal(t, uint64(0xFFC0000000000000), buildPathMask(2), "A depth 2 the mask is 0xFFC0000000000000")
	// ...
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFC), buildPathMask(28), "A depth 28 the mask is 0xFFFFFFFFFFFFFFFC")
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), buildPathMask(29), "A depth 29 the mask is 0xFFFFFFFFFFFFFFFF")
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), buildPathMask(30), "A depth 30 the mask is 0xFFFFFFFFFFFFFFFF")
}

func TestAtDepth( t *testing.T ) {

	var geo = NewGeoTri()
	var geoAt0 = geo.AtDepth(0).(geoTriData)
	var geoAt1 = geo.AtDepth(1).(geoTriData)
	var geoAt29 = geo.AtDepth(29).(geoTriData)

	assert.Equal(t, geo, geoAt0, "The origin stay the same at depth 0.")
	assert.Equal(t, geo, geoAt1, "The origin stay the same at depth 1.")
	assert.Equal(t, geo, geoAt29, "The origin stay the same at depth 29.")

	geo = NewGeoTri(1, 0)
	geoAt0 = geo.AtDepth(0).(geoTriData)
	geoAt1 = geo.AtDepth(1).(geoTriData)
	geoAt29 = geo.AtDepth(29).(geoTriData)

	assert.Equal(t, []GeoTile {1}, geoAt0.GetPath(), "The path at depth 0 is 1")
	assert.Equal(t, geo, geoAt1, "The geo triangle is defined at depth 1.")
	/*assert.Equal(t, geoAt1, geoAt29, "The geo triangle is the same at depth 1 or 29.")

	geo = geoTriData{ depth: 29, code: uint64(0x1234567890ABCDEF) }
	geoAt0 = geo.AtDepth(0).(geoTriData)
	geoAt1 = geo.AtDepth(1).(geoTriData)
	geoAt29 = geo.AtDepth(29).(geoTriData)

	assert.Equal(t, geoAt0.code, uint64(0x1000000000000000), "The code at depth 0 is 0x1000...")
	assert.Equal(t, geoAt1.code, uint64(0x1200000000000000), "The code at depth 1 is 0x1200...")
	assert.Equal(t, geo, geoAt29, "The geo triangle is defined at depth 29")
	*/
}

/*func TestGetCodeAt( t *testing.T ) {

	var geo = geoTriData{ depth: 29, code: uint64(0x1234567890ABCDEF) }

	var code, err = geo.getCodeAt(0)
	assert.Nil(t, err)
	assert.Equal(t, 4, code, "The root face is 4")

	code, err = geo.getCodeAt(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, code, "The face at depth 1 is 2")

	code, err = geo.getCodeAt(2)
	assert.Nil(t, err)
	assert.Equal(t, 0, code, "The face at depth 2 is 0")

	code, err = geo.getCodeAt(MAX_DEPTH)
	assert.Nil(t, err)
	assert.Equal(t, 3, code, "The face at depth MAX_DEPTH is 3")

	_, err = geo.getCodeAt(MAX_DEPTH + 1)
	assert.Equal(t, ERR_INVALID_DEPTH, err)
}

func TestFindNCA( t *testing.T ) {

	var geo, nca GeoTri
	var err error

	// geo = GeoTri{ depth: 6, code: uint64(0x012C000000000000) }
	// nca, err = geo.findNCA( EAST )
	
	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 2, code: uint64(0x0100000000000000) }, nca, "The nca is at depth 2 and the code is 0x0100000000000000")

	// geo = GeoTri{ depth: 1, code: uint64(0x0000000000000000) }
	// nca, err = geo.findNCA( NORTH )

	// assert.Equal(t, ERR_REACHED_SPHERE, err, "The sphere is the nca.")
	// assert.Equal(t, GeoTri{ depth: 0, code: uint64(0x0000000000000000) }, nca, "The nca is at depth 0 and the code is 0x0000000000000000")

	// geo = GeoTri{ depth: 1, code: uint64(0x4000000000000000) }
	// nca, err = geo.findNCA( NORTH )

	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 0, code: uint64(0x4000000000000000) }, nca, "The nca is at depth 0 and the code is 0x4000000000000000")

	geo = GeoTri{ depth: 1, code: uint64(0x3300000000000000) }
	nca, err = geo.findNCA( NORTH )

	assert.Equal(t, ERR_REACHED_SPHERE, err, "The sphere is the nca.")
	assert.Equal(t, GeoTri{ depth: 0, code: uint64(0x3000000000000000) }, nca, "The nca is at depth 0 and the code is 0x3000000000000000")

}

func TestFindNextSibling( t *testing.T ) {

	var depth_nca uint8
	var geo, next GeoTri
	var err error

	// depth_nca = uint8(2)
	// geo = GeoTri{ depth: 6, code: uint64(0x012C000000000000) }
	// next, err = geo.findNextSibling(EAST, depth_nca + 1)

	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 6, code: uint64(0x013C000000000000) }, next, "The next sibling is at depth 6 and the code is 0x013C000000000000")

	// depth_nca = uint8(0)
	// geo = GeoTri{ depth: 1, code: uint64(0x4000000000000000) }
	// next, err = geo.findNextSibling(EAST, depth_nca)

	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x4400000000000000) }, next, "The next sibling is at depth 1 and the code is 0x4400000000000000")

	// depth_nca = uint8(0)
	// geo = GeoTri{ depth: 1, code: uint64(0x0000000000000000) }
	// next, err = geo.findNextSibling(NORTH, depth_nca)

	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x3C00000000000000) }, next, "The next sibling is at depth 1 and the code is 0x3C00000000000000")

	// depth_nca = uint8(0)
	// geo = GeoTri{ depth: 1, code: uint64(0x4000000000000000) }
	// next, err = geo.findNextSibling(NORTH, depth_nca + 1)

	// assert.Nil(t, err)
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x4200000000000000) }, next, "The next sibling is at depth 1 and the code is 0x4200000000000000")

	depth_nca = uint8(0)
	geo = GeoTri{ depth: 1, code: uint64(0x3300000000000000) }
	next, err = geo.findNextSibling(NORTH, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x2300000000000000) }, next, "The next sibling is at depth 1 and the code is 0x2300000000000000")

}

func TestFollowToNeighbor( t *testing.T ) {
	
	depth_nca := uint8(2)
	next := GeoTri{ depth: 6, code: uint64(0x013C000000000000) }
	neighbor, err := next.followToNeighbor(EAST, depth_nca + 1)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 6, code: uint64(0x0131400000000000) }, neighbor, "The neighbor is at depth 6 and the code is 0x0131400000000000")

	depth_nca = uint8(0)
	next = GeoTri{ depth: 1, code: uint64(0x4400000000000000) }
	neighbor, err = next.followToNeighbor(EAST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x4400000000000000) }, neighbor, "The neighbor is at depth 1 and the code is 0x4400000000000000")

	depth_nca = uint8(0)
	next = GeoTri{ depth: 1, code: uint64(0x2300000000000000) }
	neighbor, err = next.followToNeighbor(EAST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x2000000000000000) }, neighbor, "The neighbor is at depth 1 and the code is 0x2000000000000000")

	depth_nca = uint8(0)
	next = GeoTri{ depth: 1, code: uint64(0x3300000000000000) }
	neighbor, err = next.followToNeighbor(WEST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x3200000000000000) }, neighbor, "The neighbor is at depth 1 and the code is 0x3200000000000000")

	depth_nca = uint8(0)
	next = GeoTri{ depth: 1, code: uint64(0x3C00000000000000) }
	neighbor, err = next.followToNeighbor(NORTH, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(0x3C00000000000000) }, neighbor, "The neighbor is at depth 1 and the code is 0x3C00000000000000")

}

func TestFindNeighbor( t *testing.T ) {

	//   EAST, NORTH, WEST, SOUTH
	var rootNeighbors = [][]uint8 {
		{0x04, 0x3C,  0x10, 0x02}, // 0x00
		{0x44, 0x42,  0x3C, 0x48}, // 0x40
		{0x20, 0x21,  0x32, 0x47}, // 0x33
		{0x32, 0x1F,  0x1C, 0x45}, // 0x31
		{0x17, 0x02,  0x15, 0x14}, // 0x16
		{0x15, 0x01,  0x27, 0x3A}, // 0x38
	}
	var geo GeoTri

//	geo = GeoTri{ depth: 1, code: uint64(0x0000000000000000) }
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[0][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x0000000000000000 -> EAST -> 0x0400000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[0][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x0000000000000000 -> NORTH -> 0x3C00000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[0][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x0000000000000000 -> WEST -> 0x1000000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[0][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x0000000000000000 -> SOUTH -> 0x0200000000000000")

//	geo = GeoTri{ depth: 1, code: uint64(0x4000000000000000) }
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[1][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x4000000000000000 -> EAST -> 0x4400000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[1][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x4000000000000000 -> NORTH -> 0x4200000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[1][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x4000000000000000 -> WEST -> 0x3C00000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[1][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x4000000000000000 -> SOUTH -> 0x4800000000000000")

	geo = GeoTri{ depth: 1, code: uint64(0x3300000000000000) }
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[2][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x3300000000000000 -> EAST -> 0x2000000000000000")
	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[2][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x3300000000000000 -> NORTH -> 0x2100000000000000")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[2][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x3300000000000000 -> WEST -> 0x3200000000000000")
	//assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[2][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3300000000000000 -> SOUTH -> 0x4700000000000000")
	
	// geo = GeoTri{ depth: 1, code: uint64(0x3100000000000000) }
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x3100000000000000 -> EAST -> 0x3200000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x3100000000000000 -> NORTH -> 0x1F00000000000000")
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x3100000000000000 -> WEST -> 0x1C00000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3100000000000000 -> SOUTH -> 0x4500000000000000")
	
	// geo = GeoTri{ depth: 1, code: uint64(0x1600000000000000) }
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[4][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x1600000000000000 -> EAST -> 0x1700000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[4][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x1600000000000000 -> NORTH -> 0x0200000000000000")
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[4][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x1600000000000000 -> WEST -> 0x1500000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[4][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x1600000000000000 -> SOUTH -> 0x1400000000000000")
	
	// geo = GeoTri{ depth: 1, code: uint64(0x3800000000000000) }
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][EAST]) << 56 }, geo.FindNeighbor(EAST), "0x3800000000000000 -> EAST -> 0x1500000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x3800000000000000 -> NORTH -> 0x0100000000000000")
	// assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][WEST]) << 56 }, geo.FindNeighbor(WEST), "0x3800000000000000 -> WEST -> 0x2700000000000000")
	// //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3800000000000000 -> SOUTH -> 0x3A00000000000000")
	
}*/