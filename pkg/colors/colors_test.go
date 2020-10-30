package colors

import (
	"fmt"
	"image/color"
	"testing"
)

var refClr = &color.RGBA{142, 12, 242, 255}

func TestFromHex(t *testing.T) {
	if clr, err := FromHex("8e0cf2"); err != nil {
		t.Error("failed parsing hex:", err)
	} else if !rgbaEquals(clr, refClr) {
		t.Errorf("color is unequal ref color: %+v", clr)
	}

	if clr, err := FromHex("#8e0cf2"); err != nil {
		t.Error("failed parsing hex:", err)
	} else if !rgbaEquals(clr, refClr) {
		t.Errorf("color is unequal ref color: %+v", clr)
	}

	if clr, err := FromHex("#8e0cf2ff"); err != nil {
		t.Error("failed parsing hex:", err)
	} else if !rgbaEquals(clr, refClr) {
		t.Errorf("color is unequal ref color: %+v", clr)
	}

	if _, err := FromHex(""); err == nil {
		t.Error("no error returned on empty string")
	}

	if _, err := FromHex("zzzzzz"); err == nil {
		t.Error("no error returned on invalid hex val")
	}
}

func TestToInt(t *testing.T) {
	if i := ToInt(refClr); i != 9309426 {
		t.Errorf("result color number was %d", i)
	}
}

func TestToHex(t *testing.T) {
	if h := ToHex(refClr); h != "8E0CF2" {
		t.Errorf("result color hex was %s", h)
	}

	if h := ToHex(&color.RGBA{0, 0, 0, 255}); h != "000000" {
		t.Errorf("result color hex was %s", h)
	}
}

func TestCreateImage(t *testing.T) {
	const (
		sizeX = 16
		sizeY = 8
	)

	imgHexData := "89504e470d0a1a0a0000000d49484452000000100000000808020000007f14e8c00000001a49444154789c62e9e3f9c4400a602249f5a8062201200000ffff6d62019f6525f13f0000000049454e44ae426082"

	if _, err := CreateImage(refClr, 0, sizeY); err == nil {
		t.Error("no error when sizeX = 0")
	}

	if _, err := CreateImage(refClr, sizeX, 0); err == nil {
		t.Error("no error when sizeY = 0")
	}

	buff, err := CreateImage(refClr, sizeX, sizeY)
	if err != nil {
		t.Error(err)
	}

	if d := fmt.Sprintf("%x", buff.Bytes()); d != imgHexData {
		t.Errorf(
			"wrong image data:\n"+
				"was:  %s\n"+
				"must: %s\n", d, imgHexData)
	}
}

func rgbaEquals(c1, c2 *color.RGBA) bool {
	return c1.R == c2.R &&
		c1.G == c2.G &&
		c1.B == c2.B &&
		c1.A == c2.A
}
