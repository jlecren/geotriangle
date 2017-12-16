package geotriangle

import (
	"fmt"
	"strings"
)

/**
 *	GeoTriError : Error type for GeoTriangle package
 */

type GeoTriError string

func (e GeoTriError) Error() string { return string(e) }

const (
	MAX_DEPTH uint8 = 29

	ERR_INVALID_DEPTH = GeoTriError("Invalid depth. Max depth is 29.")
)

/**
 *	GeoDirection : defines a direction on the sphere
 */

type GeoDirection uint8

const (
	EAST GeoDirection = iota
	NORTH
	WEST
	SOUTH
)

var geoDirectionNames = [...]string {
	"EAST",
	"NORTH",
	"WEST",
	"SOUTH",
}

func (d GeoDirection) String() string {
	return geoDirectionNames[d]
}

/**
 *	GeoTile : defines a tile type
 */

type GeoTile uint8

const (
	VERT GeoTile = iota
	LEFT
	CENTER
	RIGHT
	NONE = 0xFF
)

var geoTileNames = [...]string {
	"VERT",
	"LEFT",
	"CENTER",
	"RIGHT",
}

func (t GeoTile) String() string {
	return geoTileNames[t]
}

/**
 *	GeoTri : defines a coordinate of a triangle on a sphere
 */

type GeoTri interface {
	AtDepth( depth uint8 ) GeoTri
	FindNeighbor( dir GeoDirection ) GeoTri
	GetDepth() uint8
	GetPath() []GeoTile
	GetTileAt( depth uint8 ) (GeoTile, error)
}

func NewGeoTri(path ...GeoTile) GeoTri {
	n := uint8(len(path))
	if(n == 0) {
		return geoTriData{depth: uint8(NONE), code: uint64(0)}
	}
	if(n > MAX_DEPTH + 1) {
		n = MAX_DEPTH + 1
	}
	var code uint64
	for i := uint8(0); i < n; i++ {
		code = code << 2 + uint64(path[i])
	}
	code <<= (2 * (MAX_DEPTH - n + 1))
	return geoTriData{depth: uint8(n - 1), code: code}
}

type geoTriData struct {
	depth uint8
	code  uint64
}

func (g geoTriData) AtDepth( depth uint8 ) GeoTri {
	if( depth >= g.depth || g.depth == NONE) {
		return g
	}
	return geoTriData{ depth: depth, code: g.code & buildPathMask( depth ) }
}

func (g geoTriData) FindNeighbor( dir GeoDirection ) GeoTri {
	println("geo: ", g.String())
	var nca, err = g.findNCA(dir)
	var depth = nca.GetDepth()
	fmt.Println("NCA: ", nca)
	if(depth != NONE) {
		depth++
	} else {
		depth = 0
		println(" The NCA is the sphere.")
	}
	var next geoTriData
	next, err = g.findNextSibling(dir, depth)
	fmt.Println("Next: ", next)
	if(err != nil) {
		println(" err: ", err.Error())
		return g
	}
	next, err = next.followToNeighbor(dir, depth)
	fmt.Println("Neighbor: ", next)
	if(err != nil) {
		println(" err: ", err.Error())
		return g
	}
	return next
}

func (g geoTriData) GetDepth() uint8 {
	return g.depth
}

func (g geoTriData) GetPath() []GeoTile {
	depth := g.GetDepth()
	if(depth == NONE) {
		return make([]GeoTile, 0)
	}
	path := make([]GeoTile, g.GetDepth() + 1)
	path[0], _ = g.GetTileAt(0)
	for i := uint8(1); i < uint8(len(path)); i++ {
		path[i], _ = g.GetTileAt(i)
	}
	return path
}

func (g geoTriData) GetTileAt( depth uint8 ) (GeoTile, error) {
	if( depth == 0 ) {
		return GeoTile((g.code >> (2 * MAX_DEPTH)) & 0x3F), nil
	}
	if( depth > MAX_DEPTH ) {
		return NONE, ERR_INVALID_DEPTH
	}
	var shift = MAX_DEPTH - depth
	return GeoTile((g.code >> (2 * shift)) & 0x03), nil
}

func buildPathMask( depth uint8 ) uint64 {
	if(depth >= MAX_DEPTH) {
		return ^uint64(0)
	}
	return ^uint64(0) << (2 * (MAX_DEPTH - depth))
}

func (g geoTriData) String() string {
	depth := g.GetDepth()
	pathname := "SPHERE"
	if(depth != NONE) {
		path := g.GetPath()
		pathStr := make([]string, len(path))
		pathStr[0] = fmt.Sprintf("%02d", path[0])
		for i := uint8(1); i < uint8(len(path)); i++ {
			pathStr[i] = fmt.Sprintf("%02b", path[i])
		}
		pathname = strings.Join(pathStr, ".")
	}
	return fmt.Sprintf("{ depth: %2d, code: %0#16X, path: %s }", depth, g.code, pathname)
}

func (g *geoTriData) setCodeAt( depth uint8, code GeoTile ) error {
	println(fmt.Sprintf("setCodeAt - depth: %d code: %06b", depth, code))
	if( depth > MAX_DEPTH ) {
		return ERR_INVALID_DEPTH
	}
	var newCode uint64 = uint64(code)
	var mask uint64
	if( depth == 0 ) {
		mask = uint64(0x3F) << (2 * MAX_DEPTH)
		newCode <<= (2 * MAX_DEPTH)
	} else {
		mask = uint64(0x03) << (2 * (MAX_DEPTH - depth))
		newCode <<= (2 * (MAX_DEPTH - depth))
	}
	println(fmt.Sprintf("before: %0#16x", g.code))
	println(fmt.Sprintf("   add: %0#16x", newCode))
	println(fmt.Sprintf("  mask: %0#16x", mask))
	println(fmt.Sprintf(" ^mask: %0#16x", ^mask))
	g.code = ( g.code & ^mask ) + (newCode & mask)
	println(fmt.Sprintf("after : %0#16x", g.code))
	return nil
}

  // EAST,  NORTH, WEST
var stopTab = [][]bool{
    {false, false,  false}, // VERT
	{true,  false, false}, // LEFT
	{true,  true,  true},  // CENTER
	{false, false, true},  // RIGHT
}

/**
 * Find the Nearest Common Ancestor
 */
func (g geoTriData) findNCA( dir GeoDirection ) (geoTriData, error) {
	println(fmt.Sprintf("findNCA - dir: %s", dir.String()))
	var depth = g.GetDepth()
	var child_type, err = g.GetTileAt(depth)
	if( err != nil ) {
		return g, err
	}
	println(fmt.Sprintf("depth: %d code: %06b", depth, child_type))
	for( depth > 0 && !stopTab[child_type][dir]) {
		depth--
		child_type, err = g.GetTileAt(depth)
		if( err != nil ) {
			return g, err
		}
	}
	if( depth == 0) {
		topRow := child_type / 5
		println(fmt.Sprintf("topRow: %d", topRow))
		if( ((topRow % 2 == 1) && (dir == NORTH)) ||
			((topRow % 2 == 0) && (dir == SOUTH))) {
			return g.AtDepth(0).(geoTriData), nil
		}
		return NewGeoTri().(geoTriData), nil
	}
	return g.AtDepth(depth - 1).(geoTriData), nil
}

//   EAST,   NORTH,  WEST
var siblingTab = [][]GeoTile{
    {LEFT,   CENTER, RIGHT},  // VERT
	{CENTER, RIGHT,   VERT},   // LEFT
	{RIGHT,  VERT,   LEFT},   // CENTER
	{VERT,   LEFT,  CENTER}, // RIGHT
}

//   EAST,   NORTH,  WEST, SOUTH
var siblingTop = [][]GeoTile{
    {GeoTile(1), 3, 4, 5}, // 0
	{GeoTile(2), 4, 0, 6}, // 1
	{GeoTile(3), 0, 1, 7}, // 2
	{GeoTile(4), 1, 2, 8}, // 3
	{GeoTile(0), 2, 3, 9}, // 4
	{GeoTile(10), 0, 14, 10}, // 5
	{GeoTile(11), 1, 10, 11}, // 6
	{GeoTile(12), 2, 11, 12}, // 7
	{GeoTile(13), 3, 12, 13}, // 8
	{GeoTile(14), 4, 13, 14}, // 9
	{GeoTile(6), 5, 5, 15}, // 10
	{GeoTile(7), 6, 6, 16}, // 11
	{GeoTile(8), 7, 7, 17}, // 12
	{GeoTile(9), 8, 8, 18}, // 13
	{GeoTile(5), 9, 9, 19}, // 14
	{GeoTile(16), 10, 19, 0}, // 15
	{GeoTile(17), 11, 15, 1}, // 16
	{GeoTile(18), 12, 16, 2}, // 17
	{GeoTile(19), 13, 17, 3}, // 18
	{GeoTile(15), 14, 18, 4}, // 19
}

func (g geoTriData) findNextSibling( dir GeoDirection, depth_nca uint8 ) (geoTriData, error) {
	var tile_type, err = g.GetTileAt(depth_nca)
	if(err != nil) {
		return g, err
	}
	var next_tile GeoTile
	if( depth_nca > 0 ) {
		next_tile = siblingTab[tile_type][dir]
	} else {
		next_tile = siblingTop[tile_type][dir]
		if(10 <= tile_type && tile_type <= 14 && g.GetDepth() > 0) {
			subtile, _ := g.GetTileAt(depth_nca + 1)
			if(subtile == 3) {
				next_tile++
				if(next_tile == 14) {
					next_tile = 5
				}
			}
		}
	}
	err = g.setCodeAt(depth_nca, GeoTile(next_tile))
	if(err != nil) {
		return g, err
	} 
	return g, nil
}

//           EAST,  NORTH, WEST
var reflTop = [][]GeoTile {
    {GeoTile(VERT), VERT,  VERT},  // VERT
	{GeoTile(NONE), LEFT,  RIGHT}, // LEFT
	{GeoTile(NONE), CENTER,  NONE},  // CENTER
	{GeoTile(LEFT), RIGHT, NONE},  // RIGHT
}

func (g geoTriData) followToNeighbor( dir GeoDirection, depth_nca uint8 ) (geoTriData, error) {
	println(fmt.Sprintf("followToNeighbor - dir: %s depth_nca: %d", dir.String(), depth_nca))
	depth := g.GetDepth()
	tile_type, err := g.GetTileAt(depth_nca)
	println(fmt.Sprintf("geo: %s", g.String()))
	println(fmt.Sprintf("NCA tile: %d %06b", tile_type, tile_type))
	if( depth_nca > 0 || ( 5 <= tile_type && tile_type <= 14)) {
		for ( depth_nca < depth ) {
			depth_nca++
			tile_type, err = g.GetTileAt(depth_nca)
			println(fmt.Sprintf("tile at %d:    %02b", depth_nca, tile_type))
			println(fmt.Sprintf("sibling at %d: %02b", depth_nca, uint8(siblingTab[tile_type][dir])))
			if(err != nil) {
				return g, err
			}
			err = g.setCodeAt(depth_nca, siblingTab[tile_type][dir])
			if(err != nil) {
				return g, err
			}
		}
	} else {
		for ( depth_nca < depth ) {
			depth_nca++
			tile_type, err = g.GetTileAt(depth_nca)
			if(err != nil) {
				return g, err
			}
			var newCode = reflTop[tile_type][dir]
			if( newCode == 0xFF ) {
				continue
			}
			err = g.setCodeAt(depth_nca, newCode)
			if(err != nil) {
				return g, err
			}
		}
	}

	return g, nil
}