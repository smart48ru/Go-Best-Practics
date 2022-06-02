//go:build integration
// +build integration

package integration

import (
	"Best-GO/internal/scann"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func Test_Scanner(t *testing.T) {
	sc := scann.New(10, "../../test", ".txt", 2)
	go sc.FindFiles()

	//
	fl := []string{}
	// map для записи в нее всех файлов считанных из дирикторий для дальнейшего сравнения
	sm := make(map[string]struct{})
	// slice для записи в него имен файлов не содержащиеся в мапе
	//exam := []string{}
	tests := []struct {
		name    string
		want    []string
		exam    []string
		sm      map[string]struct{}
		wantErr bool
	}{
		{"All file OK",
			[]string{"../../test/1/1-1.txt", "../../test/1/1-2.txt", "../../test/2/2-1.txt", "../../test/2/2-2.txt", "../../test/2/2-3.txt"},
			[]string{},
			make(map[string]struct{}),
			false,
		},
		{"File /1/1-3.txt not in dir",
			[]string{"../../test/1/1-3.txt", "../../test/1/1-1.txt", "../../test/1/1-2.txt", "../../test/2/2-1.txt", "../../test/2/2-2.txt", "../../test/2/2-3.txt"},
			[]string{},
			make(map[string]struct{}),
			true,
		},
		{"File /1/1-1.txt not in wanted",
			[]string{"../../test/1/1-2.txt", "../../test/2/2-1.txt", "../../test/2/2-2.txt", "../../test/2/2-3.txt"},
			[]string{},
			make(map[string]struct{}),
			true,
		},
	}

	for {
		select {
		case <-sc.Ctx().Done():
			log.Info().Msg("Done")
			for _, n := range fl {
				if _, ok := sm[n]; !ok { // проверка на всякий случай. если вдруг кто-то уже записал данное значение или фаловая система поломалась
					sm[n] = struct{}{}
				}
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					for _, n := range tt.want {
						if _, ok := sm[n]; !ok {
							tt.exam = append(tt.exam, n)
						}
					}
					if len(tt.exam) != 0 && !tt.wantErr {
						for _, f := range tt.exam {
							t.Errorf("file = %v", f)
						}
					}
				})
			}

			sc.CtxCancel()
			return
		case err := <-sc.ErrChan():
			assert.Nil(t, err, "Error in scanner")
		case result := <-sc.ResChan():
			// Комментарии пишу для себя
			// записываем в слайс все имена файлов которые прочитали из директорий
			fl = append(fl, string(result.Path()))
		}
	}

}
