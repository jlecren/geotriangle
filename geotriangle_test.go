package geotriangle

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

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

	var (
		geo = NewGeoTri()
		geoAt0 = geo.AtDepth(0).(geoTriData)
		geoAt1 = geo.AtDepth(1).(geoTriData)
		geoAt29 = geo.AtDepth(29).(geoTriData)
	)

	assert.Equal(t, geo, geoAt0, "The origin stay the same at depth 0.")
	assert.Equal(t, geo, geoAt1, "The origin stay the same at depth 1.")
	assert.Equal(t, geo, geoAt29, "The origin stay the same at depth 29.")

	geo = NewGeoTri(1, 0)
	geoAt0 = geo.AtDepth(0).(geoTriData)
	geoAt1 = geo.AtDepth(1).(geoTriData)
	geoAt29 = geo.AtDepth(29).(geoTriData)

	assert.Equal(t, []GeoTile {1}, geoAt0.GetPath(), "The path at depth 0 is 1")
	assert.Equal(t, geo, geoAt1, "The geo triangle is defined at depth 1.")
	assert.Equal(t, geoAt1, geoAt29, "The geo triangle is the same at depth 1 or 29.")

	geo = geoTriData{ depth: 29, code: uint64(0x1234567890ABCDEF) }
	geoAt0 = geo.AtDepth(0).(geoTriData)
	geoAt1 = geo.AtDepth(1).(geoTriData)
	geoAt29 = geo.AtDepth(29).(geoTriData)

	assert.Equal(t, geoAt0.code, uint64(0x1000000000000000), "The code at depth 0 is 0x1000...")
	assert.Equal(t, geoAt1.code, uint64(0x1200000000000000), "The code at depth 1 is 0x1200...")
	assert.Equal(t, geo, geoAt29, "The geo triangle is defined at depth 29")
	
}

func TestGetTileAt( t *testing.T ) {

	var (
		geo = NewGeoTri(
			4, 2, 0, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 3).(geoTriData)
		code, err = geo.GetTileAt(0)
	)

	assert.Nil(t, err)
	assert.Equal(t, GeoTile(4), code, "The root face is 4")

	code, err = geo.GetTileAt(1)
	assert.Nil(t, err)
	assert.Equal(t, GeoTile(2), code, "The face at depth 1 is 2")

	code, err = geo.GetTileAt(2)
	assert.Nil(t, err)
	assert.Equal(t, GeoTile(0), code, "The face at depth 2 is 0")

	code, err = geo.GetTileAt(MAX_DEPTH)
	fmt.Println("geo: ", geo)
	assert.Nil(t, err)
	assert.Equal(t, GeoTile(3), code, "The face at depth MAX_DEPTH is 3")

	_, err = geo.GetTileAt(MAX_DEPTH + 1)
	assert.Equal(t, ERR_INVALID_DEPTH, err)
}

func TestFindNCA( t *testing.T ) {

	var (
		geo, nca geoTriData
		err error
	)

	geo = NewGeoTri(0, 1, 0, 2, 3, 0, 0).(geoTriData)
	nca, err = geo.findNCA( EAST )
	
	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 1, 0), nca, "The nca is at depth 2 and the code is 0x0100000000000000")

	geo = NewGeoTri(0, 0).(geoTriData)
	nca, err = geo.findNCA( NORTH )

	assert.Equal(t, NewGeoTri(), nca, "The nca is at depth FF and the code is 0x0000000000000000")

	geo = NewGeoTri(16, 0).(geoTriData)
	nca, err = geo.findNCA( NORTH )

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(16), nca, "The nca is at depth 0 and the code is 0x4000000000000000")

	geo = NewGeoTri(12, 3).(geoTriData)
	nca, err = geo.findNCA( NORTH )

	assert.Equal(t, NewGeoTri(), nca, "The nca is at depth FF and the code is 0x3000000000000000")

}

func TestFindNextSibling( t *testing.T ) {

	var (
		depth_nca uint8
		geo, next geoTriData
		err error
	)

	depth_nca = uint8(2)
	geo = NewGeoTri(0, 1, 0, 2, 3, 0, 0).(geoTriData)
	next, err = geo.findNextSibling(EAST, depth_nca + 1)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 1, 0, 3, 3, 0, 0), next, "The next sibling is at depth 6 and the code is 0x013C000000000000")

	depth_nca = uint8(0)
	geo = NewGeoTri(16, 0).(geoTriData)
	next, err = geo.findNextSibling(EAST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(17, 0), next, "The next sibling is at depth 1 and the code is 0x4400000000000000")

	depth_nca = uint8(0)
	geo = NewGeoTri(0, 0).(geoTriData)
	next, err = geo.findNextSibling(NORTH, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(15, 0), next, "The next sibling is at depth 1 and the code is 0x3C00000000000000")

	depth_nca = uint8(0)
	geo = NewGeoTri(16, 0).(geoTriData)
	next, err = geo.findNextSibling(NORTH, depth_nca + 1)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(16, 2), next, "The next sibling is at depth 1 and the code is 0x4200000000000000")

	depth_nca = uint8(0)
	geo = NewGeoTri(12, 3).(geoTriData)
	next, err = geo.findNextSibling(NORTH, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(8, 3), next, "The next sibling is at depth 1 and the code is 0x2300000000000000")

}

func TestFollowToNeighbor( t *testing.T ) {
	
	depth_nca := uint8(2)
	next := NewGeoTri(0, 1, 0, 3, 3, 0, 0).(geoTriData)
	neighbor, err := next.followToNeighbor(EAST, depth_nca + 1)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 1, 0, 3, 0, 1, 1), neighbor, "The neighbor is at depth 6 and the code is 0x0131400000000000")

	depth_nca = uint8(0)
	next = NewGeoTri(1, 1).(geoTriData)
	neighbor, err = next.followToNeighbor(EAST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(1, 1), neighbor, "The neighbor is at depth 1 and the code is 0x4400000000000000")

	depth_nca = uint8(0)
	next = NewGeoTri(0, 2).(geoTriData)
	neighbor, err = next.followToNeighbor(EAST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 2), neighbor, "The neighbor is at depth 1 and the code is 0x2000000000000000")

	depth_nca = uint8(0)
	next = NewGeoTri(0, 3).(geoTriData)
	neighbor, err = next.followToNeighbor(WEST, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 3), neighbor, "The neighbor is at depth 1 and the code is 0x3200000000000000")

	depth_nca = uint8(0)
	next = NewGeoTri(0, 3).(geoTriData)
	neighbor, err = next.followToNeighbor(NORTH, depth_nca)

	assert.Nil(t, err)
	assert.Equal(t, NewGeoTri(0, 3), neighbor, "The neighbor is at depth 1 and the code is 0x3C00000000000000")

}

func TestFindNeighbor( t *testing.T ) {
	var geo GeoTri

	geo = NewGeoTri(0, 0)
	assert.Equal(t, NewGeoTri(1, 0), geo.FindNeighbor(EAST), "(0, 0) -> EAST -> (1, 0)")
	assert.Equal(t, NewGeoTri(3, 0), geo.FindNeighbor(NORTH), "(0, 0) -> NORTH -> (3, 0)")
	assert.Equal(t, NewGeoTri(4, 0), geo.FindNeighbor(WEST), "(0, 0) -> WEST -> (4, 0)")
	// assert.Equal(t, NewGeoTri(0, 2), geo.FindNeighbor(SOUTH), "(0, 0) -> SOUTH -> (0, 2)")

	geo = NewGeoTri(1, 0)
	assert.Equal(t, NewGeoTri(2, 0), geo.FindNeighbor(EAST), "(1, 0) -> EAST -> (2, 0)")
	assert.Equal(t, NewGeoTri(4, 0), geo.FindNeighbor(NORTH), "(1, 0) -> NORTH -> (4, 0)")
	assert.Equal(t, NewGeoTri(0, 0), geo.FindNeighbor(WEST), "(1, 0) -> WEST -> (0, 0)")
//	assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[1][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x4000000000000000 -> SOUTH -> 0x4800000000000000")

	geo = NewGeoTri(12, 3)
	assert.Equal(t, NewGeoTri(8, 0), geo.FindNeighbor(EAST), "(12, 3) -> EAST -> (8, 0)")
	assert.Equal(t, NewGeoTri(8, 1), geo.FindNeighbor(NORTH), "(12, 3) -> NORTH -> (8, 1)")
	// assert.Equal(t, NewGeoTri(12, 2), geo.FindNeighbor(WEST), "(12, 3) -> WEST -> (12, 2)")
	//assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[2][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3300000000000000 -> SOUTH -> 0x4700000000000000")
	
	// geo = NewGeoTri(12, 1)
	// assert.Equal(t, NewGeoTri(12, 2), geo.FindNeighbor(EAST), "(12, 1) -> EAST -> (12, 2)")
	// // //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x3100000000000000 -> NORTH -> 0x1F00000000000000")
	// assert.Equal(t, NewGeoTri(7, 0), geo.FindNeighbor(WEST), "(12, 1) -> WEST -> (7, 0)")
	// // //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[3][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3100000000000000 -> SOUTH -> 0x4500000000000000")
	
	// geo = NewGeoTri(5, 2)
	// assert.Equal(t, NewGeoTri(5, 3), geo.FindNeighbor(EAST), "(5, 2) -> EAST -> (5, 3)")
	// // assert.Equal(t, NewGeoTri(0, 2), geo.FindNeighbor(NORTH), "(5, 2) -> NORTH -> (0, 2)")
	// assert.Equal(t, NewGeoTri(5, 1), geo.FindNeighbor(WEST), "(5, 2) -> WEST -> (5, 1)")
	// // assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[4][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x1600000000000000 -> SOUTH -> 0x1400000000000000")
	
	// geo = NewGeoTri(15, 0)
	// assert.Equal(t, NewGeoTri(16, 0), geo.FindNeighbor(EAST), "(15, 0) -> EAST -> (16, 0)")
	// // //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][NORTH]) << 56 }, geo.FindNeighbor(NORTH), "0x3800000000000000 -> NORTH -> 0x0100000000000000")
	// assert.Equal(t, NewGeoTri(19, 0), geo.FindNeighbor(WEST), "(15, 0) -> WEST -> (19, 0)")
	// // //assert.Equal(t, GeoTri{ depth: 1, code: uint64(rootNeighbors[5][SOUTH]) << 56 }, geo.FindNeighbor(SOUTH), "0x3800000000000000 -> SOUTH -> 0x3A00000000000000")
	
}