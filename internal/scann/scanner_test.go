package scann

import (
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func NewScan() (sc scanner, dir string, ext string, depth int64, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return scanner{}, "", "", 0, err
	}
	ext = ".go"
	depth = 2

	sc = New(10, dir, ext, depth)
	return sc, dir, ext, depth, nil
}

func TestNew(t *testing.T) {
	c, dir, ext, depth, err := NewScan()
	if err != nil {
		log.Error().Msgf("TestNew %s", err)
	}
	assert.NotEmpty(t, c.ctx, "new scanner ctx is Empty")
	assert.NotEmpty(t, c.ctxCancel, "new scanner ctxCancel is Empty")
	assert.Empty(t, c.wg, "new scanner wg is NOT Empty")
	assert.Empty(t, c.errChan, "new scanner errChan is NOT Empty")
	assert.Empty(t, c.resChan, "new scanner resChan is NOT Empty")
	assert.NotEmpty(t, c.dir, "new scanner dir is Empty")
	assert.Empty(t, c.cDir, "new scanner cDir is NOT Empty")
	assert.Equal(t, c.dir, dir)
	assert.NotEmpty(t, c.ext, "new scanner ctx is Empty")
	assert.Equal(t, c.ext, ext)
	assert.NotEmpty(t, c.depth, "new scanner ctx is Empty")
	assert.Equal(t, c.depth, depth)
}

func TestScanner_DeIncDepth(t *testing.T) {
	c, _, _, _, err := NewScan()
	if err != nil {
		log.Error().Msgf("Test_DeIncDepth %s", err)
	}
	assert.Equal(t, c.Depth(), int64(2))
	c.DeIncDepth()
	assert.Equal(t, c.Depth(), int64(1))
	c.DeIncDepth()
	assert.Equal(t, c.Depth(), int64(0))
}

func TestScanner_IncDepth(t *testing.T) {
	c, _, _, _, err := NewScan()
	if err != nil {
		log.Error().Msgf("Test_IncDepth %s", err)
	}
	assert.Equal(t, c.Depth(), int64(2))
	c.IncDepth()
	assert.Equal(t, c.Depth(), int64(4))
	c.IncDepth()
	assert.Equal(t, c.Depth(), int64(6))
}

//Test_Scanner_Dir - test Read dir
func Test_Scanner_Dir(t *testing.T) {
	sc := New(10, "/11/11", ".txt", 2)
	go sc.FindFiles()
	for {
		select {
		case <-sc.Ctx().Done():
			log.Info().Msg("Done")
			sc.CtxCancel()
			return
		case err := <-sc.ErrChan():
			assert.NotNil(t, err, "Error in scanner")
		}
	}
}

//Test_Scanner_Depth - Test depth
func Test_Scanner_Depth(t *testing.T) {
	sc := New(10, "../../test", ".txt", -1)
	go sc.FindFiles()

	for {
		select {
		case <-sc.Ctx().Done():
			log.Info().Msg("Done")
			sc.CtxCancel()
			return
		case err := <-sc.ErrChan():
			assert.Nil(t, err, "Error in scanner")
		case result := <-sc.ResChan():
			assert.Nil(t, result, "Error in scanner  - Depth not work")
		}
	}
}
