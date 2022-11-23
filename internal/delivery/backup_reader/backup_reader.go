package backup_reader

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/vbua/go_setter_getter/internal/entity"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserGradeSetManyService interface {
	SetMany(map[string]entity.UserGrade)
}

type BackupReader struct {
	userGradeSetManyService UserGradeSetManyService
}

func NewBackupReader(userGradeSetManyService UserGradeSetManyService) BackupReader {
	return BackupReader{userGradeSetManyService}
}

func (b BackupReader) ReadFromBackup() time.Time {
	resp, err := http.Get(os.Getenv("NEXT_REPLICA_URL") + "/backup")
	if err == nil {
		defer resp.Body.Close()
		_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
		if err != nil {
			log.Fatalln(err)
		}
		sepFilename := strings.Split(params["filename"], ".")
		if len(sepFilename) == 0 {
			return time.Time{}
		}
		parseTime, err := time.Parse("01-02-2006 15:04:05", sepFilename[0])

		if err != nil {
			log.Fatalln(err)
		}

		gr, err := gzip.NewReader(resp.Body)
		if err == io.EOF { // значит пустой файл
			return time.Time{}
		}
		if err != nil { // если формат не тот, значит что-то в целом неправильно работает в приложениях
			log.Fatalln(err)
		}
		defer gr.Close()
		cr := csv.NewReader(gr)
		rec, err := cr.ReadAll()
		if err != nil {
			log.Fatalln(err)
		}
		userGrades := make(map[string]entity.UserGrade)
		for _, grade := range rec {
			postpaidLimit, err := strconv.Atoi(grade[1])
			if err != nil {
				fmt.Println(err)
			}
			spp, err := strconv.Atoi(grade[2])
			if err != nil {
				fmt.Println(err)
			}
			shippingFee, err := strconv.Atoi(grade[3])
			if err != nil {
				fmt.Println(err)
			}
			returnFee, err := strconv.Atoi(grade[4])
			if err != nil {
				fmt.Println(err)
			}
			userGrades[grade[0]] = entity.UserGrade{
				UserId:        grade[0],
				PostpaidLimit: postpaidLimit,
				Spp:           spp,
				ShippingFee:   shippingFee,
				ReturnFee:     returnFee,
			}
		}
		if len(userGrades) > 0 {
			b.userGradeSetManyService.SetMany(userGrades)
		}

		return parseTime
	}

	return time.Time{}
}
