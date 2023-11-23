// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package shaping

import (
	"unicode"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/harfbuzz"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/loader"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/unicode/bidi"
)

type Input struct {
	// Text is the body of text being shaped. Only the range Text[RunStart:RunEnd] is considered
	// for shaping, with the rest provided as context for the shaper. This helps with, for example,
	// cross-run Arabic shaping or handling combining marks at the start of a run.
	Text []rune
	// RunStart and RunEnd indicate the subslice of Text being shaped.
	RunStart, RunEnd int
	// Direction is the directionality of the text.
	Direction di.Direction
	// Face is the font face to render the text in.
	Face font.Face

	// FontFeatures activates or deactivates optional features
	// provided by the font.
	// The settings are applied to the whole [Text].
	FontFeatures []FontFeature

	// Size is the requested size of the font.
	// More generally, it is a scale factor applied to the resulting metrics.
	// For instance, given a device resolution (in dpi) and a point size (like 14), the `Size` to
	// get result in pixels is given by : pointSize * dpi / 72
	Size fixed.Int26_6

	// Script is an identifier for the writing system used in the text.
	Script language.Script

	// Language is an identifier for the language of the text.
	Language language.Language
}

// FontFeature sets one font feature.
//
// A font feature is an optionnal behavior a font might expose,
// identified by a 4 bytes [Tag].
// Most features are disabled by default; setting a non zero [Value]
// enables it.
//
// An exemple of font feature is the replacement of fractions (like 1/2, 3/4)
// by specialized glyphs, which would be activated by using
//
//	FontFeature{Tag: loader.MustNewTag("frac"), Value: 1}
//
// See also https://learn.microsoft.com/en-us/typography/opentype/spec/featurelist
// and https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_fonts/OpenType_fonts_guide
type FontFeature struct {
	Tag   loader.Tag
	Value uint32
}

// Fontmap provides a general mechanism to select
// a face to use when shaping text.
type Fontmap interface {
	// ResolveFace is called by `SplitByFace` for each input rune potentially
	// triggering a face change.
	// It must always return a valid (non nil) font.Face value.
	ResolveFace(r rune) font.Face
}

var _ Fontmap = fixedFontmap(nil)

type fixedFontmap []font.Face

// ResolveFace panics if the slice is empty
func (ff fixedFontmap) ResolveFace(r rune) font.Face {
	for _, f := range ff {
		if _, has := f.NominalGlyph(r); has {
			return f
		}
	}
	return ff[0]
}

// SplitByFontGlyphs split the runes from 'input' to several items, sharing the same
// characteristics as 'input', expected for the `Face` which is set to
// the first font among 'availableFonts' providing support for all the runes
// in the item.
// Runes supported by no fonts are mapped to the first element of 'availableFonts', which
// must not be empty.
// The 'Face' field of 'input' is ignored: only 'availableFaces' are consulted.
// Rune coverage is obtained by calling the NominalGlyph() method of each font.
// See also SplitByFace for a more general approach of font selection.
func SplitByFontGlyphs(input Input, availableFaces []font.Face) []Input {
	return SplitByFace(input, fixedFontmap(availableFaces))
}

// SplitByFace split the runes from 'input' to several items, sharing the same
// characteristics as 'input', expected for the `Face` which is set to
// the return value of the `Fontmap.ResolveFace` call.
// The 'Face' field of 'input' is ignored: only 'availableFaces' is used to select the face.
func SplitByFace(input Input, availableFaces Fontmap) []Input {
	var splitInputs []Input
	currentInput := input
	for i := input.RunStart; i < input.RunEnd; i++ {
		r := input.Text[i]
		if currentInput.Face != nil && ignoreFaceChange(r) {
			// add the rune to the current input
			continue
		}

		// select the first font supporting r
		selectedFace := availableFaces.ResolveFace(r)

		if currentInput.Face == selectedFace {
			// add the rune to the current input
			continue
		}

		// new face needed

		if i != input.RunStart {
			// close the current input ...
			currentInput.RunEnd = i
			// ... add it to the output ...
			splitInputs = append(splitInputs, currentInput)
		}

		// ... and create a new one
		currentInput = input
		currentInput.RunStart = i
		currentInput.Face = selectedFace
	}

	// close and add the last input
	currentInput.RunEnd = input.RunEnd
	splitInputs = append(splitInputs, currentInput)
	return splitInputs
}

// ignoreFaceChange returns `true` is the given rune should not trigger
// a change of font.
//
// We don't want space characters to affect font selection; in general,
// it's always wrong to select a font just to render a space.
// We assume that all fonts have the ASCII space, and for other space
// characters if they don't, HarfBuzz will compatibility-decompose them
// to ASCII space...
//
// We don't want to change fonts for line or paragraph separators.
//
// Finaly, we also don't change fonts for what Harfbuzz consider
// as ignorable (however, some Control Format runes like 06DD are not ignored).
//
// The rationale is taken from pango : see bugs
// https://bugzilla.gnome.org/show_bug.cgi?id=355987
// https://bugzilla.gnome.org/show_bug.cgi?id=701652
// https://bugzilla.gnome.org/show_bug.cgi?id=781123
// for more details.
func ignoreFaceChange(r rune) bool {
	return unicode.Is(unicode.Cc, r) || // control
		unicode.Is(unicode.Cs, r) || // surrogate
		unicode.Is(unicode.Zl, r) || // line separator
		unicode.Is(unicode.Zp, r) || // paragraph separator
		(unicode.Is(unicode.Zs, r) && r != '\u1680') || // space separator != OGHAM SPACE MARK
		harfbuzz.IsDefaultIgnorable(r)
}

// Segmenter holds a state used to split input
// according to three caracteristics : text direction (bidi),
// script, and face.
type Segmenter struct {
	// pools of inputs, used to reduce allocations,
	// which are alternatively considered as input/output of the segmentation
	buffer1, buffer2 []Input

	// buffer used for bidi segmentation
	bidiParagraph bidi.Paragraph

	// used to handle Common script
	parenthesisStack []int
}

// Split segments the given [text] according to :
//   - text direction
//   - script
//   - face, as defined by [faces]
//
// As a consequence, the following fields of the returned runs are set :
//   - Text, RunStart, RunEnd
//   - Direction
//   - Script
//   - Face
//
// [defaultDirection] is used during bidi ordering, and should refer to the general
// context [text] is used in (typically the user system preference for GUI apps.)
//
// The returned sliced is owned by the [Segmenter] and is only valid until
// the next call to [Split].
func (seg *Segmenter) Split(text []rune, faces Fontmap, defaultDirection di.Direction) []Input {
	seg.reset()
	seg.splitByBidi(text, defaultDirection)
	seg.splitByScript()
	seg.splitByFace()
	return seg.buffer1
}

func (seg *Segmenter) reset() {
	// zero the slices to avoid 'memory leak' on pointer slice fields
	for i := range seg.buffer1 {
		seg.buffer1[i].Text = nil
		seg.buffer1[i].FontFeatures = nil
	}
	for i := range seg.buffer2 {
		seg.buffer2[i].Text = nil
		seg.buffer2[i].FontFeatures = nil
	}
	seg.buffer1 = seg.buffer1[:0]
	seg.buffer2 = seg.buffer2[:0]
	// bidiParagraph is reset when using SetString
	seg.parenthesisStack = seg.parenthesisStack[:0]
}

// fills buffer1
func (seg *Segmenter) splitByBidi(text []rune, defaultDirection di.Direction) {
	if defaultDirection.Axis() != di.Horizontal || len(text) == 0 {
		seg.buffer1 = append(seg.buffer1, Input{
			Text:      text,
			RunStart:  0,
			RunEnd:    len(text),
			Direction: defaultDirection,
		})
		return
	}
	def := bidi.LeftToRight
	if defaultDirection.Progression() == di.TowardTopLeft {
		def = bidi.RightToLeft
	}
	seg.bidiParagraph.SetString(string(text), bidi.DefaultDirection(def))
	out, err := seg.bidiParagraph.Order()
	if err != nil {
		seg.buffer1 = append(seg.buffer1, Input{
			Text:      text,
			RunStart:  0,
			RunEnd:    len(text),
			Direction: defaultDirection,
		})
		return
	}

	input := Input{Text: text} // start a rune 0
	for i := 0; i < out.NumRuns(); i++ {
		currentInput := input
		run := out.Run(i)
		dir := run.Direction()
		_, endRune := run.Pos()
		currentInput.RunEnd = endRune + 1
		if dir == bidi.RightToLeft {
			currentInput.Direction = di.DirectionRTL
		} else {
			currentInput.Direction = di.DirectionLTR
		}
		seg.buffer1 = append(seg.buffer1, currentInput)
		input.RunStart = currentInput.RunEnd
	}
}

// uses buffer1 as input and fills buffer2
func (seg *Segmenter) splitByScript() {}

// uses buffer2 as input, resets and fills buffer1
func (seg *Segmenter) splitByFace() {}
