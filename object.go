package tilepix

import (
	"fmt"

	"github.com/faiface/pixel"

	log "github.com/sirupsen/logrus"
)

/*
   ___  _     _        _
  / _ \| |__ (_)___ __| |_
 | (_) | '_ \| / -_) _|  _|
  \___/|_.__// \___\__|\__|
           |__/
*/

// Object is a TMX file struture holding a specific Tiled object.
type Object struct {
	Name       string      `xml:"name,attr"`
	Type       string      `xml:"type,attr"`
	X          float64     `xml:"x,attr"`
	Y          float64     `xml:"y,attr"`
	Width      float64     `xml:"width,attr"`
	Height     float64     `xml:"height,attr"`
	GID        int         `xml:"id,attr"`
	Visible    bool        `xml:"visible,attr"`
	Polygon    *Polygon    `xml:"polygon"`
	PolyLine   *PolyLine   `xml:"polyline"`
	Properties []*Property `xml:"properties>property"`
	Ellipse    *struct{}   `xml:"ellipse"`
	Point      *struct{}   `xml:"point"`

	objectType ObjectType

	// parentMap is the map which contains this object
	parentMap *Map
}

// GetRect will return a pixel.Circle representation of this object relative to the map (the co-ordinates will match
// those as drawn in Tiled).  If the object type is not `EllipseObj` this function will return `pixel.C(pixel.ZV, 0)`
// and an error.
//
// Because there is no pixel geometry code for irregular ellipses, this function will average the width and height of
// the ellipse object from the TMX file, and return a regular circle about the centre of the ellipse.
func (o *Object) GetEllipse() (pixel.Circle, error) {
	if o.GetType() != EllipseObj {
		log.WithError(ErrInvalidObjectType).WithField("Object type", o.GetType()).Error("Object.GetEllipse: object type mismatch")
		return pixel.C(pixel.ZV, 0), ErrInvalidObjectType
	}

	// In TMX files, ellipses are defined by the containing rectangle.  The X, Y positions are the bottom-left (after we
	// have flipped them).
	// Because Pixel does not support irregular ellipses, we take the average of width and height.
	radius := (o.Width + o.Height) / 4
	// The centre should be the same as the ellipses drawn in Tiled, this will make outputs more intuitive.
	centre := pixel.V(o.X+(o.Width/2), o.Y+(o.Height/2))

	return pixel.C(centre, radius), nil
}

// GetPoint will return a pixel.Vec representation of this object relative to the map (the co-ordinates will match those
// as drawn in Tiled).  If the object type is not `PointObj` this function will return `pixel.ZV` and an error.
func (o *Object) GetPoint() (pixel.Vec, error) {
	if o.GetType() != PointObj {
		log.WithError(ErrInvalidObjectType).WithField("Object type", o.GetType()).Error("Object.GetPoint: object type mismatch")
		return pixel.ZV, ErrInvalidObjectType
	}

	return pixel.V(o.X, o.Y), nil
}

// GetRect will return a pixel.Rect representation of this object relative to the map (the co-ordinates will match those
// as drawn in Tiled).  If the object type is not `RectangleObj` this function will return `pixel.R(0, 0, 0, 0)` and an
// error.
func (o *Object) GetRect() (pixel.Rect, error) {
	if o.GetType() != RectangleObj {
		log.WithError(ErrInvalidObjectType).WithField("Object type", o.GetType()).Error("Object.GetRect: object type mismatch")
		return pixel.R(0, 0, 0, 0), ErrInvalidObjectType
	}

	return pixel.R(o.X, o.Y, o.X+o.Width, o.Y+o.Height), nil
}

// GetType will return the ObjectType constant type of this object.
func (o *Object) GetType() ObjectType {
	return o.objectType
}

func (o *Object) String() string {
	return fmt.Sprintf("Object{%s, Name: '%s'}", o.objectType, o.Name)
}

func (o *Object) flipY() {
	o.Y = o.parentMap.pixelHeight() - o.Y - o.Height
}

// hydrateType will work out what type this object is.
func (o *Object) hydrateType() {
	if o.Polygon != nil {
		o.objectType = PolygonObj
		return
	}

	if o.PolyLine != nil {
		o.objectType = PolylineObj
		return
	}

	if o.Ellipse != nil {
		o.objectType = EllipseObj
		return
	}

	if o.Point != nil {
		o.objectType = PointObj
		return
	}

	o.objectType = RectangleObj
}

func (o *Object) setParent(m *Map) {
	o.parentMap = m

	if o.Polygon != nil {
		o.Polygon.setParent(m)
	}
	if o.PolyLine != nil {
		o.PolyLine.setParent(m)
	}
	for _, p := range o.Properties {
		p.setParent(m)
	}
}
